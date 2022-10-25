package bot

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"strings"

	db_action "github.com/0ne-zero/BookShopBot/internal/database/action"
	"github.com/0ne-zero/BookShopBot/internal/database/model"
	"github.com/0ne-zero/BookShopBot/internal/utils"
	setting "github.com/0ne-zero/BookShopBot/internal/utils/settings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	exif_r "github.com/scottleedavis/go-exif-remove"
	persian_time "github.com/yaa110/go-persian-calendar"
)

type validateUserinputFunc func(*tgbotapi.Update) error

func ConfigBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	bot, err := tgbotapi.NewBotAPI(API_KEY)
	if err != nil {
		return nil, nil, fmt.Errorf("error occurred during create new bot instance - %w", err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	updates := bot.GetUpdatesChan(u)

	cfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "/start", Description: "start the bot"},
		tgbotapi.BotCommand{Command: "/test", Description: "test the bot"},
	)
	bot.Request(cfg)
	return bot, updates, nil
}

func SendUnknownError(bot *tgbotapi.BotAPI, chat_id int64) {
	_, err := bot.Send(tgbotapi.NewMessage(chat_id, UNKNOWN_ERROR))
	if err != nil {
		log.Printf("Error occurred during send error message - %s\n", err.Error())
	}
}
func SendError(bot *tgbotapi.BotAPI, chat_id int64, error_text string) {
	_, err := bot.Send(tgbotapi.NewMessage(chat_id, error_text))
	if err != nil {
		log.Printf("Error occurred during send error message - %s\n", err.Error())
	}
}
func IsCommand(text string) bool {
	return strings.HasPrefix(text, "/")
}
func makeBookCoverTypesKeyboard() (*tgbotapi.InlineKeyboardMarkup, error) {
	types, err := db_action.GetBookCoverTypes()
	if err != nil {
		return nil, err
	}
	var keyboard_items []tgbotapi.InlineKeyboardButton
	for i := range types {
		item := tgbotapi.NewInlineKeyboardButtonData(types[i].Type, fmt.Sprint(types[i].ID))
		keyboard_items = append(keyboard_items, item)
	}
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboard_items...),
	)
	return &keyboard, nil
}
func makeBookSizeKeyboard() (*tgbotapi.InlineKeyboardMarkup, error) {
	sizes, err := db_action.GetBookSize()
	if err != nil {
		return nil, err
	}
	var keyboard_items []tgbotapi.InlineKeyboardButton
	for i := range sizes {
		item := tgbotapi.NewInlineKeyboardButtonData(sizes[i].Name, fmt.Sprint(sizes[i].ID))
		keyboard_items = append(keyboard_items, item)
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	keyboard_items_len := len(keyboard_items)
	if keyboard_items_len < 4 {
		rows = append(rows, keyboard_items)
	} else {
		start := 0
		end := 3
		number_of_rows_divide := float64(keyboard_items_len) / float64(3)
		if utils.IsFloatNumberRound(number_of_rows_divide) {
			for i := 0; i < int(number_of_rows_divide); i++ {
				rows = append(rows, keyboard_items[start:end])
				start += 3
				end += 3
			}
		} else {
			number_of_rows_divide++
			for i := 0; i < int(number_of_rows_divide); i++ {
				if end > keyboard_items_len {
					rows = append(rows, keyboard_items[start:])
				}
				rows = append(rows, keyboard_items[start:end])
				start += 3
				end += 3
			}
		}
	}
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard, nil
}
func makeBookKeyboard(book_id int) *tgbotapi.InlineKeyboardMarkup {
	callback_data := fmt.Sprint(ADD_BOOK_TO_CART_INLINE_KEYBOARD_ITEM_TITLE, "?", book_id)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(ADD_BOOK_TO_CART_INLINE_KEYBOARD_ITEM_TITLE, callback_data)))
	return &keyboard
}
func makeBookExistsInCartKeyboard(book_id int) *tgbotapi.InlineKeyboardMarkup {
	callback_data := fmt.Sprint(DELETE_BOOK_FROM_CART_INLINE_KEYBOARD_ITEM_TITLE, "?", book_id)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(DELETE_BOOK_FROM_CART_INLINE_KEYBOARD_ITEM_TITLE, callback_data)))
	return &keyboard
}
func makeContactToAdminInlineKeyboard() (*tgbotapi.InlineKeyboardMarkup, error) {
	admin_username := setting.ReadFieldInSettingData("ADMIN_TELEGRAM_USERNAME")
	if admin_username == "" {
		log.Printf("Error occurred during read ADMIN_TELEGRAM_USERNAME field from setting file - Field is empty")
		return nil, fmt.Errorf("error occurred during read ADMIN_TELEGRAM_USERNAME field from setting file - Field is empty")
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(CONTACT_ADMIN_KEYBOARD_ITEM_TITLE, fmt.Sprintf(SWITCH_TO_PV_FORMAT, admin_username))))
	return &keyboard, nil
}
func extractBookIDFromCallbackData(data string) string {
	return data[strings.LastIndex(data, "?")+1:]
}
func getInputFromUser(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel, input_request_text string, validate_func validateUserinputFunc) (string, error) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, input_request_text)
	msg.ReplyMarkup = BACK_TO_MAIN_MENU_KEYBOARD
	_, err := bot_api.Send(msg)
	if err != nil {
		log.Printf("Error occurred during send request for book title - %s", err.Error())
		return "", err
	}
	// Wait for user input
	for inner_update := range *updates {
		if err = validate_func(&inner_update); err != nil {
			return "", err
		}
		return inner_update.Message.Text, nil
	}
	return "", fmt.Errorf("nothing happened")
}
func IsStartQuery(text string) bool {
	if strings.Contains(text, "https://t.me/Xbookshopbot/?start=") || strings.Contains(text, "/start") && text != "/start" {
		return true
	} else {
		return false
	}
}
func extractBookIDFromStartQuery(query string) (int, error) {
	splitted := strings.Split(query, "=")
	var id string
	// If query doesn't have "=" character, It's diffrent start query mode
	if len(splitted) == 1 {
		splitted = strings.Split(query, " ")
		if len(splitted) == 1 {
			return 0, fmt.Errorf("query doesn't have id or we cannot extract id")
		}
		id = splitted[1]
	} else {
		id = splitted[1]
	}

	return strconv.Atoi(id)
}
func formatBookInformation(book_id int) (string, error) {
	book, err := db_action.GetBookByID(book_id)
	if err != nil {
		return "", err
	}
	var censor_status string
	if book.Censored {
		censor_status = "سانسور شده"
	} else {
		censor_status = "بدون سانسور"
	}
	formatted_info := fmt.Sprintf(
		BOOK_INFORMATION_FORMAT, book.Title, book.Author, book.Translator, book.NumberOfPages, book.Genre,
		censor_status, book.CoverType.Type, book.BookSize.Name, book.BookAgeCategory.Category, fmt.Sprint(book.GoodReadsScore),
		fmt.Sprint(book.ArezoScore), book.Publisher, book.PublishDate, book.ISBN, fmt.Sprint(book.Price), BOT_USERNAME)
	return formatted_info, nil
}
func GetGoodreadsScoreByISBN(isbn string) (string, error) {
	isbn = strings.ReplaceAll(isbn, "-", "")
	cmd := exec.Command("python3", "../internal/urils/isbn_to_goodreads_score.py", isbn)
	buffer := bytes.Buffer{}
	cmd.Stdout = &buffer
	output_bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output_bytes), nil
}
func IsAdmin(user_telegram_id int) bool {
	// Get admin id from settings
	admin_id_str := setting.ReadFieldInSettingData("ADMIN_TELEGRAM_ID")
	// Parse id to int64
	admin_id, err := strconv.ParseInt(admin_id_str, 10, 64)
	if err != nil {
		return false
	}
	// Check ids are equal or not
	if int64(user_telegram_id) == admin_id {
		return true
	} else {
		return false
	}
}
func removeExifFromPhoto(bytes []byte) ([]byte, error) {
	removed_exif_bytes, err := exif_r.Remove(bytes)
	if err != nil {
		return nil, fmt.Errorf("error occurred during removing exif from file")
	}
	return removed_exif_bytes, nil
}
func generateRandomBytes(size int) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := rand.Read(bytes)
	return bytes, err
}
func generateRandomHex(length int) (string, error) {
	var byte_size = length
	if length%2 != 0 {
		byte_size += 1
	}
	bytes, err := generateRandomBytes(byte_size / 2)
	if err != nil {
		return "", err
	}
	hex := hex.EncodeToString(bytes)
	hex_len := len(hex)
	for hex_len != length {
		hex = hex[:hex_len-1]
		hex_len = len(hex)
	}
	return hex, nil
}

// Returns saved photo path
func downloadAndSavePhoto(download_url string) (string, error) {
	// Download photo
	pic_bytes, err := utils.DownloadFileFromURL(download_url)
	if err != nil {
		log.Printf("Error occurred during downloading document/file - %s\n", err.Error())
		return "", err
	}

	// Remove exif metadata
	pic_bytes, err = removeExifFromPhoto(pic_bytes)
	if err != nil {
		log.Printf("Error occurred during remove exif data from sent photo - %s", err.Error())
		return "", err
	}

	// Generate random pic name
	pic_name, err := generateRandomHex(32)
	if err != nil {
		log.Printf("Error occurred during generate random file name - %s", err.Error())
		return "", err
	}

	// Create file path
	pic_path := filepath.Join(PICTURES_DIRECTORY, pic_name, ".jpg")

	// Save picture in local disk
	err = utils.WriteBytesToFile(pic_path, pic_bytes)
	if err != nil {
		log.Printf("Error occurred during save picture in localdisk - %s", err.Error())
		return "", err
	}
	return pic_path, nil
}
func getUserAddressInformationFromUser(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) (*model.Address, error) {
	var addr model.Address
	var err error
	var input_fetched = false
	// Get address country
	for !input_fetched {
		addr.Country, err = getInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_COUNTRY, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address country - %s", err.Error())
			return nil, err
		}
		input_fetched = true
	}

	// Get address province
	input_fetched = false
	for !input_fetched {
		addr.Province, err = getInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_PROVINCE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address province - %s", err.Error())
			return nil, err
		}
		input_fetched = true
	}
	// Get address city
	input_fetched = false
	for !input_fetched {
		addr.City, err = getInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_CITY, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address city - %s", err.Error())
			return nil, err
		}
		input_fetched = true
	}
	// Get address street
	input_fetched = false
	for !input_fetched {
		addr.Street, err = getInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_STREET, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address street - %s", err.Error())
			return nil, err
		}
		input_fetched = true
	}
	// Get address building number
	input_fetched = false
	for !input_fetched {
		addr.BuildingNumber, err = getInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_BUILDING_NUMBER, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address building number - %s", err.Error())
			return nil, err
		}
		input_fetched = true
	}
	// Get address postal code
	input_fetched = false
	for !input_fetched {
		addr.PostalCode, err = getInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_POSTAL_CODE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address postal code - %s", err.Error())
			return nil, err
		}
		input_fetched = true
	}
	// Get address description
	input_fetched = false
	for !input_fetched {
		addr.Description, err = getInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_DESCRIPTION, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address description - %s", err.Error())
			return nil, err
		}
		input_fetched = true
	}
	// Get address phone number
	input_fetched = false
	for !input_fetched {
		addr.PhoneNumber, err = getInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_PHONE_NUMBER, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address phone number - %s", err.Error())
			return nil, err
		}
		input_fetched = true
	}
	return &addr, nil
}
func makeCartMessage(user_telegram_id int) (string, error) {
	var message string
	// Add message header
	message += CART_MESSAGE_HEADER + "\n\n"
	message = "محتویات سبد خرید:\n"
	// Get user cart books id
	books_id, err := db_action.GetUserCartBooksID(user_telegram_id)
	if err != nil {
		log.Printf("Error occurred during get user cart books id - %s", err.Error())
		return "", nil
	}
	// Show contents of cart
	for i := range books_id {
		// Get book name
		book_name, err := db_action.GetBookTitleByID(books_id[i])
		if err != nil {
			log.Printf("Error occurred during get book title by book id - %s", err.Error())
			return "", err
		}
		// Get book author
		book_author, err := db_action.GetBookAuthorByID(books_id[i])
		if err != nil {
			log.Printf("Error occurred during get book author by book id - %s", err.Error())
			return "", err
		}
		item := fmt.Sprintf("%d- <a href=\"%s\">%s (%s)\n</a>", i+1, fmt.Sprintf(BOT_START_QUERY, BOT_USERNAME, books_id[i]), book_name, book_author)
		// Append item to message
		message += item
	}
	// Append price of cart (without shipment cost)
	// Calculate books price
	books_price, err := db_action.GetUserCartTotalPrice(user_telegram_id)
	if err != nil {
		return "", nil
	}
	message += "\n" + fmt.Sprintf("قیمت کتاب ها:  %s\n", fmt.Sprint(books_price))
	// Create message footer
	message += "\n" + CART_MESSAGE_FOOTER
	// Add bot username at the end of message
	message += fmt.Sprintf("\n\n@%s", BOT_USERNAME)
	return message, nil
}
func makeBuyCartMessage(user_telegram_id int) (string, error) {
	// Get store cart number from setting file
	store_cart_number := setting.ReadFieldInSettingData("STORE_CART_NUMBER")
	if store_cart_number == "" {
		return "", fmt.Errorf("STORE_CART_NUMBER is empty, so program can't perform its porpuse")
	}
	// Get store cart number owner full name
	store_cart_number_owner_fullname := setting.ReadFieldInSettingData("STORE_CART_NUMBER_OWNER_FULLNAME")
	if store_cart_number_owner_fullname == "" {
		return "", fmt.Errorf("STORE_CART_NUMBER_OWNER_FULLNAME is empty")
	}
	// Craete message and append header to it
	var message string = BUY_CART_MESSAGE_HEADER_MESSAGE
	// Append card information to message
	message += fmt.Sprintf("\n\n%s\n%s\n", fmt.Sprintf("شماره کارت فروشگاه: %s", store_cart_number), fmt.Sprintf("مالک شماره کارت: %s", store_cart_number_owner_fullname))

	// Calculate cart total price (books + shipment cost)
	books_price, shipment_cost, err := calculateCartTotalPrice(user_telegram_id)
	if err != nil {
		return "", err
	}
	// Append total price to message
	message += "\n" + fmt.Sprintf("قیمت کتاب ها: : %s\nهزینه ی ارسال  : %s\n<b>قیمت نهایی</b>: %s\n", fmt.Sprint(books_price), fmt.Sprint(shipment_cost), fmt.Sprint(books_price+shipment_cost))
	// Append footer to message
	message += BUY_CART_MESSAGE_FOOTER_MESSAGE
	// Add bot username at the end of message
	message += fmt.Sprintf("\n\n@%s", BOT_USERNAME)
	return message, nil
}
func makeShowUserOrdersMessage(user_telegram_id int) (string, error) {
	// Create message and append header of message
	var message string = SHOW_ORDERS_HEADER_MESSAGE
	orders_info, err := db_action.GetUserOrdersForShowByUserTelegramID(user_telegram_id)
	if err != nil {
		return "", nil
	}
	// Add orders info to message
	for i := range orders_info {
		cur_order := orders_info[i]
		// Add order time and status
		message += fmt.Sprintf("%d- %s (%s) *%s*", i+1, "جزيیات سفارش ثبت شده در تاریخ :", ConvertTimeToPersian(cur_order.OrderTime), cur_order.OrderStatus)
		// Add order books to message
		for i := range cur_order.Books {
			message += fmt.Sprintf("\n\t%s(%s)", cur_order.Books[i].Title, cur_order.Books[i].Author)
		}
		message += "\n\n"
	}
	// Append footer of message
	message += SHOW_ORDERS_FOOTER_MESSAGE
	return message, nil
}

// Calculate price of cart (books + shipment cost)
func calculateCartTotalPrice(user_telegram_id int) (float32, float32, error) {
	// Calculate books price
	cart_price, err := db_action.GetUserCartTotalPrice(user_telegram_id)
	if err != nil {
		return 0, 0, nil
	}
	// Get cart books and address needed information for calculate shipment cost
	cart_info, err := db_action.GetCartInformationForCalculateShipmentCost(user_telegram_id)
	if err != nil {
		return 0, 0, err
	}
	// Calculate cart shipment cost
	shipment_cost, err := calculateShipmentCost(cart_info)
	if err != nil {
		return 0, 0, err
	}
	// Cart total price
	return cart_price, shipment_cost, nil
}
func calculateShipmentCost(cart_info *db_action.CartInformationForCalculateShipmentCost) (float32, error) {
	return 0, nil
}
func ConvertTimeToPersian(t *time.Time) string {
	p_time := persian_time.New(*t)
	format := fmt.Sprintf("%s-%s-%s", p_time.Year(), p_time.Month(), p_time.Day())
	return format
}
func makeMainKeyboard(user_telegram_id int) (*tgbotapi.ReplyKeyboardMarkup, error) {
	// Check user is admin
	if IsAdmin(user_telegram_id) {
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(ADMIN_STATISTICS_KEYBOARD_ITEM_TITLE),
				tgbotapi.NewKeyboardButton(ADMIN_CONFIRM_ORDERS_KEYBOARD_ITEM_TITLE),
				tgbotapi.NewKeyboardButton(ADMIN_DELETE_BOOK_KEYBOARD_ITEM_TITLE),
				tgbotapi.NewKeyboardButton(ADMIN_ADD_BOOK_KEYBOARD_ITEM_TITLE),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(FAQ_KEYBOARD_ITEM_TITLE),
				tgbotapi.NewKeyboardButton(ADMIN_BACK_TO_USER_PANEL_ITEM_TITLE)),
		)
		return &keyboard, nil
		// User isn't admin
	} else {
		if have_order, err := db_action.DoesUserHaveOrder(user_telegram_id); err != nil {
			log.Printf("Error occurred during check does user have any order - %s", err.Error())
			return nil, err
			// User have order
		} else if have_order {
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(CART_KEYBOARD_ITEM_TITLE),
					tgbotapi.NewKeyboardButton(ADDRESS_KEYBOARD_ITEM_TITLE),
					tgbotapi.NewKeyboardButton(SEARCH_BOOK_KEYBOARD_ITEM_TITLE),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(SHOW_ORDERS_KEYBOARD_ITEM_TITLE),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(FAQ_KEYBOARD_ITEM_TITLE),
					tgbotapi.NewKeyboardButton(CONTACT_ADMIN_KEYBOARD_ITEM_TITLE)),
			)
			return &keyboard, nil
			// User doesn't have order
		} else {
			keyboard := tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(CART_KEYBOARD_ITEM_TITLE),
					tgbotapi.NewKeyboardButton(ADDRESS_KEYBOARD_ITEM_TITLE),
					tgbotapi.NewKeyboardButton(SEARCH_BOOK_KEYBOARD_ITEM_TITLE),
				),
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton(FAQ_KEYBOARD_ITEM_TITLE),
					tgbotapi.NewKeyboardButton(CONTACT_ADMIN_KEYBOARD_ITEM_TITLE)),
			)
			return &keyboard, nil
		}
	}
}
func makeAdminUserPanelKeyboard() *tgbotapi.ReplyKeyboardMarkup {
	var keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CART_KEYBOARD_ITEM_TITLE),
			tgbotapi.NewKeyboardButton(ADDRESS_KEYBOARD_ITEM_TITLE),
			tgbotapi.NewKeyboardButton(SEARCH_BOOK_KEYBOARD_ITEM_TITLE),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(CONTACT_ADMIN_KEYBOARD_ITEM_TITLE),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(ADMIN_BACK_TO_ADMIN_PANEL_ITEM_TITLE),
		),
	)
	return &keyboard
}
func makeBookAgeCategoryKeyboard() (*tgbotapi.InlineKeyboardMarkup, error) {
	categories, err := db_action.GetBookAgeCategories()
	if err != nil {
		return nil, err
	}
	var keyboard_items []tgbotapi.InlineKeyboardButton
	for i := range categories {
		item := tgbotapi.NewInlineKeyboardButtonData(categories[i].Category, fmt.Sprint(categories[i].ID))
		keyboard_items = append(keyboard_items, item)
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	keyboard_items_len := len(keyboard_items)
	if keyboard_items_len < 4 {
		rows = append(rows, keyboard_items)
	} else {
		start := 0
		end := 3
		number_of_rows_divide := float64(keyboard_items_len) / float64(3)
		if utils.IsFloatNumberRound(number_of_rows_divide) {
			for i := 0; i < int(number_of_rows_divide); i++ {
				rows = append(rows, keyboard_items[start:end])
				start += 3
				end += 3
			}
		} else {
			number_of_rows_divide++
			for i := 0; i < int(number_of_rows_divide); i++ {
				if end > keyboard_items_len {
					rows = append(rows, keyboard_items[start:])
				}
				rows = append(rows, keyboard_items[start:end])
				start += 3
				end += 3
			}
		}
	}
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard, nil
}
func extractMessageIDFromTelegramRawResponse(raw_response string) (int, error) {
	_, after, found := strings.Cut(raw_response, "\"message_id\":")
	if !found {
		return 0, fmt.Errorf("message_id field not found")
	}
	end_id_index := strings.Index(after, ",")
	id_str := after[:end_id_index]
	return strconv.Atoi(id_str)
}
