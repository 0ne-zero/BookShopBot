package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	db_action "github.com/0ne-zero/BookShopBot/internal/database/action"
	"github.com/0ne-zero/BookShopBot/internal/database/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var updateValidateFunc validateUserinputFunc = func(u *tgbotapi.Update) error {
	if u.Message == nil {
		return fmt.Errorf("Update.Message is nil")
	}
	if u.Message.Text == "" {
		return fmt.Errorf("Update.Message.Text is nil")
	}
	return nil
}

func SearchBookByTitle_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, SEARCH_TEXT)
	msg.ReplyMarkup = SEARCH_BOOK_INLINE_KEYBOARD
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send search book by title message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func Cart_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Check user cart is empty
	if empty, err := db_action.IsUserCartEmptyByUserTelegramID(int(update.SentFrom().ID)); err != nil {
		log.Printf("Error occurred during check user cart is empty - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		// User cart is empty
	} else if empty {
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, CART_IS_EMPTY)
		if _, err := bot_api.Send(msg); err != nil {
			log.Printf("Error occurred during send cart is empty message - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
		// User cart isn't empty
	} else {
		// Create cart message
		message, err := makeCartMessage(int(update.SentFrom().ID))
		if err != nil {
			log.Printf("Error occurred during make buy cart message - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, message)
		msg.ReplyMarkup = BUY_CART_INLINE_KEYBOARD
		msg.ParseMode = "html"
		if _, err := bot_api.Send(msg); err != nil {
			log.Printf("Error occurred during send show cart message - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
	}

}
func ContactAdmin_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, CONTACT_TO_ADMIN_MESSAGE)
	var err error
	msg.ReplyMarkup, err = makeContactToAdminInlineKeyboard()
	if err != nil {
		log.Printf("Error occurred during make contact to admin inline keyboard - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send contact to admin message - %s", err.Error())
	}
}
func SetAddress_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) {
	// Check user already has address
	var user_have_address bool
	var err error
	if user_have_address, err = db_action.DoesUserHaveAddress(int(update.SentFrom().ID)); err != nil {
		log.Printf("Error occurred during check user have address - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
	// If user already have address show their address and ask for set again it
	if user_have_address {
		// Create edit address message
		var message string
		if !strings.HasSuffix(YOU_ALREADY_HAVE_ADDRESS, "\n") {
			message += YOU_ALREADY_HAVE_ADDRESS + "\n"
		} else {
			message += YOU_ALREADY_HAVE_ADDRESS
		}
		// Get user address
		exists_addr, err := db_action.GetUserAddressByTelegramUserID(int(update.SentFrom().ID))
		if err != nil {
			log.Printf("Error occurred during get user exist address for show - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			return
		}
		message += fmt.Sprintf(SHOW_USER_ADDRESS_FORMATTED, exists_addr.Country, exists_addr.Province, exists_addr.City, exists_addr.Street, exists_addr.BuildingNumber, exists_addr.PostalCode, exists_addr.PhoneNumber, exists_addr.Description)
		message += fmt.Sprintf("\n@%s", BOT_USERNAME)
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, message)
		msg.ReplyMarkup = FOR_EDIT_ADDRESS_INLINE_KEYBOARD
		// Send edit address message
		if _, err = bot_api.Send(msg); err != nil {
			log.Printf("Error occureed during send user address for edit message - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			return
		}
	} else {
		SetAddress_InlineKeyboardHandler(bot_api, update, updates)
	}
}
func Admin_AddBook_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) {
	// Create book, It should be fill with user inputs
	book := model.Book{}
	var err error
	// It's equal to true by default, jsut for enter to below for loop
	var input_fetched bool = false
	// Get Images/Pictures of book
	for !input_fetched {
		// Send request for book picture
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, REQUEST_BOOK_PICTURE)
		if _, err := bot_api.Send(msg); err != nil {
			log.Printf("Error occurred during send book picutes request - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
		for inner_update := range *updates {
			if inner_update.Message.Document != nil {
				// Extract document/file download url
				download_url, err := bot_api.GetFileDirectURL(update.Message.Document.FileID)
				if err != nil {
					log.Printf("Error occurred during extract document/file download url - %s\n", err.Error())
					SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
					continue
				}
				// Download image
				pic_path, err := downloadAndSavePhoto(download_url)
				if err != nil {
					log.Printf("Error occurred during download and save image - %s", err.Error())
					SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
				}
				// Add path to book pictures path
				book.Pictures += pic_path + "|"
			} else if inner_update.Message.Photo != nil {
				// Get every sent photo and remove exif data form them
				// Each sent photo has four quality, so update.Message.Photo has four item, we need only the last one (original photo)
				// Extract last item (main item)
				main_photo_file_id := inner_update.Message.Photo[len(inner_update.Message.Photo)-1].FileID
				/// Extract document/file download url
				download_url, err := bot_api.GetFileDirectURL(main_photo_file_id)
				if err != nil {
					log.Printf("Error occurred during extract document/file download url - %s\n", err.Error())
					SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
					continue
				}
				// Download image
				pic_path, err := downloadAndSavePhoto(download_url)
				if err != nil {
					log.Printf("Error occurred during download and save image - %s", err.Error())
					SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
				}
				// Add path to book pictures path
				book.Pictures += pic_path + "|"
			}
		}
	}
	// Get book's weight
	input_fetched = false
	for !input_fetched {
		weight_str, err := getInputFromUser(bot_api, update, updates, REQUEST_BOOK_WEIGHT, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book weight - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
		// Convert response to float from string
		weight_float, err := strconv.ParseFloat(weight_str, 32)
		if err != nil {
			log.Printf("Error occurred during convert string weight to float32 weight - %s", err.Error())
			SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_NON_NUMBER_VALUE_ERROR)
		}
		book.Weight = float32(weight_float)
		input_fetched = true
	}
	// Get book's title
	input_fetched = false
	for !input_fetched {
		raw_isbn, err := getInputFromUser(bot_api, update, updates, REQUEST_BOOK_ISBN, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book title - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		book.ISBN = strings.ReplaceAll(raw_isbn, "-", "")
		input_fetched = true
	}
	// Get book's title
	input_fetched = false
	for !input_fetched {
		book.Title, err = getInputFromUser(bot_api, update, updates, REQUEST_BOOK_TITLE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book title - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get book's author
	input_fetched = false
	for !input_fetched {
		book.Author, err = getInputFromUser(bot_api, update, updates, REQUEST_BOOK_AUTHOR, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get book's translator
	input_fetched = false
	for !input_fetched {
		book.Translator, err = getInputFromUser(bot_api, update, updates, REQUEST_BOOK_TRANSLATOR, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get book's paper type
	input_fetched = false
	for !input_fetched {
		book.PaperType, err = getInputFromUser(bot_api, update, updates, REQUEST_BOOK_PAPER_TYPE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get book's description
	input_fetched = false
	for !input_fetched {
		book.Description, err = getInputFromUser(bot_api, update, updates, REQUEST_BOOK_DESCRIPTION, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get book's number of pages
	input_fetched = false
	for !input_fetched {
		number_of_pages_str, err := getInputFromUser(bot_api, update, updates, REQUEST_BOOK_NUMBER_OF_PAGES, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		book.NumberOfPages, err = strconv.Atoi(number_of_pages_str)
		if err != nil {
			log.Printf("Entered non-int value - %s", err.Error())
			SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_NON_NUMBER_VALUE_ERROR)
			continue
		}
		input_fetched = true
	}
	// Get book's genre
	input_fetched = false
	for !input_fetched {
		book.Genre, err = getInputFromUser(bot_api, update, updates, REQUEST_BOOK_GENRE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get book's censored status
	input_fetched = false
	for !input_fetched {
		// Send request
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, REQUEST_BOOK_CENSORED_STATUS)
		msg.ReplyMarkup = CENSORED_STATUS_KEYBOARD
		if _, err := bot_api.Send(msg); err != nil {
			log.Printf("Error occurred during send request for book title - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		// Get response
		for inner_update := range *updates {
			if inner_update.CallbackQuery == nil || inner_update.CallbackQuery.Data == "" {
				log.Print("Invalid response from user for censor status of book")
				SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_VALUE_IS_INVALID_ERROR)
				continue
			}
			// Values meaning:
			// 0 = Isn't censored
			// 1 = Is censored
			if inner_update.CallbackQuery.Data == "1" {
				book.Censored = true
				input_fetched = true
			} else if inner_update.CallbackQuery.Data == "0" {
				book.Censored = false
				input_fetched = true
			}
			break
		}
	}
	// Get book's publisher
	input_fetched = false
	for !input_fetched {
		book.Publisher, err = getInputFromUser(bot_api, update, updates, REQUEST_BOOK_PUBLISHER, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get book's publish date
	input_fetched = false
	for !input_fetched {
		book.PublishDate, err = getInputFromUser(bot_api, update, updates, REQUEST_BOOK_PUBLISHDATE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get book's price
	input_fetched = false
	for !input_fetched {
		price_str, err := getInputFromUser(bot_api, update, updates, REQUEST_BOOK_PRICE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		price_float64, err := strconv.ParseFloat(price_str, 32)
		book.Price = float32(price_float64)
		if err != nil {
			log.Printf("Entered non-float value for price- %s", err.Error())
			SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_NON_NUMBER_VALUE_ERROR)
			continue
		}
		input_fetched = true
	}
	// Get book's arezo score
	input_fetched = false
	for !input_fetched {
		arezo_score_str, err := getInputFromUser(bot_api, update, updates, REQUEST_BOOK_AZERO_SCORE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		arezo_score_float64, err := strconv.ParseFloat(arezo_score_str, 32)
		if err != nil {
			log.Printf("Entered non-float valud for arezo score")
			SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_NON_NUMBER_VALUE_ERROR)
			continue
		}
		book.ArezoScore = float32(arezo_score_float64)
		input_fetched = true
	}
	// Get book's cover type
	input_fetched = false
	for !input_fetched {
		// Get cover types from database and create keyboard
		keyboard, err := makeBookCoverTypesKeyboard()
		if err != nil {
			log.Printf("Error occurred during make book cover type keyboard - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, REQUEST_BOOK_COVERTYPE)
		msg.ReplyMarkup = keyboard
		if _, err = bot_api.Send(msg); err != nil {
			log.Printf("Error occurred during send book cover type request message - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		// Get user response
		for inner_update := range *updates {
			if inner_update.CallbackQuery == nil || inner_update.CallbackQuery.Data == "" {
				log.Print("Invalid response from user for cover type of book")
				SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_VALUE_IS_INVALID_ERROR)
				continue
			}
			data_int, err := strconv.Atoi(inner_update.CallbackQuery.Data)
			if err != nil {
				log.Printf("Entered non-int value for selected cover type id - %s", err.Error())
				SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_NON_NUMBER_VALUE_ERROR)
				continue
			}

			if selected_type, err := db_action.GetBookCoverTypeByID(data_int); err != nil || selected_type == nil {
				log.Printf("Entered unknown cover type id")
				SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_VALUE_IS_INVALID_ERROR)
				continue
			} else {
				book.CoverType = selected_type
				book.BookCoverTypeID = selected_type.ID
				input_fetched = true
			}
			break
		}
	}
	// Get book's size
	input_fetched = false
	for !input_fetched {
		keyboard, err := makeBookSizeKeyboard()
		if err != nil {
			log.Printf("Error occurred during make book size keyboard - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, REQUEST_BOOK_SIZE)
		msg.ReplyMarkup = keyboard
		if _, err = bot_api.Send(msg); err != nil {
			log.Printf("Error occurred during send book size request - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		// Get user response
		for inner_update := range *updates {
			if inner_update.CallbackQuery == nil && inner_update.CallbackQuery.Data == "" {
				log.Print("Invalid response from user for size of book")
				SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_VALUE_IS_INVALID_ERROR)
				continue
			}
			data_int, err := strconv.Atoi(inner_update.CallbackQuery.Data)
			if err != nil {
				log.Printf("Entered non-int value for selected book size - %s", err.Error())
				SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_NON_NUMBER_VALUE_ERROR)
				continue
			}
			if selected_size, err := db_action.GetBookSizeByID(data_int); err != nil || selected_size == nil {
				log.Printf("Entered unknown book size id")
				SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_VALUE_IS_INVALID_ERROR)
				continue
			} else {
				book.BookSize = selected_size
				book.BookSizeID = selected_size.ID
				input_fetched = true
			}
			break
		}
	}
	// Get book's age category
	input_fetched = false
	for !input_fetched {
		keyboard, err := makeBookAgeCategoryKeyboard()
		if err != nil {
			log.Printf("Error occurred during make age category keyboard - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, REQUEST_BOOK_AGE_CATEGORY)
		msg.ReplyMarkup = keyboard
		if _, err := bot_api.Send(msg); err != nil {
			log.Printf("Error occurred during send book age category request - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		// Get user response
		for inner_update := range *updates {
			if inner_update.CallbackQuery == nil && inner_update.CallbackQuery.Data == "" {
				log.Print("Invalid response from user for age category of book")
				SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_VALUE_IS_INVALID_ERROR)
				continue
			}
			data_int, err := strconv.Atoi(inner_update.CallbackQuery.Data)
			if err != nil {
				log.Printf("Entered non-int value for selected age category - %s", err.Error())
				SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_NON_NUMBER_VALUE_ERROR)
				continue
			}
			if selected_age_category, err := db_action.GetBookAgeCategoryByID(data_int); err != nil || selected_age_category == nil {
				log.Printf("Entered unknown book age category")
				SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_VALUE_IS_INVALID_ERROR)
				continue
			} else {
				book.BookAgeCategory = selected_age_category
				book.BookAgeCategoryID = selected_age_category.ID
				input_fetched = true
			}
		}
		price_str, err := getInputFromUser(bot_api, update, updates, REQUEST_BOOK_PRICE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		price_float64, err := strconv.ParseFloat(price_str, 32)
		book.Price = float32(price_float64)
		if err != nil {
			log.Printf("Entered non-float value for price- %s", err.Error())
			SendError(bot_api, update.FromChat().ChatConfig().ChatID, ENTERED_VALUE_IS_INVALID_ERROR)
			continue
		}
		input_fetched = true
	}

	// Get goodreads score
	if isbn_len := len(book.ISBN); isbn_len == 13 || isbn_len == 10 {
		goodreads_score_str, err := GetGoodreadsScoreByISBN(book.ISBN)
		if err != nil {
			log.Printf("Error occurred during get googreads score - %s", err.Error())
		} else {
			var goodreads_score_float float64
			if goodreads_score_float, err = strconv.ParseFloat(goodreads_score_str, 32); err != nil {
				log.Printf("Error occurred during convert string goodreads score to float - %s", err.Error())
			} else {
				book.GoodReadsScore = float32(goodreads_score_float)
			}
		}
	}

	err = db_action.AddBook(&book)
	if err != nil {
		log.Printf("Error occurred during add book to database")
		SendError(bot_api, update.FromChat().ChatConfig().ChatID, BOOK_NOT_SAVED_IN_DATABASE_ERROR)
	}
}

func Admin_DeleteBook_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) {
	// Send search keyboard to user
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, SEARCH_FOR_DELETE_BOOK_MESSAGE_TEXT)
	msg.ReplyMarkup = SEARCH_BOOK_INLINE_KEYBOARD
	if update.Message != nil && update.Message.MessageID != 0 {
		msg.ReplyToMessageID = update.Message.MessageID
	}
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send search book for delete - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
	// Get user response
	for inner_update := range *updates {
		if inner_update.InlineQuery != nil && inner_update.InlineQuery.Query != "" && inner_update.InlineQuery.Query != DELETE_STRING {
			SearchBookByTitle_InlineQueryHandler(bot_api, &inner_update)
		} else if inner_update.Message != nil && inner_update.Message.Text != "" {
			if IsStartQuery(inner_update.Message.Text) {
				book_id, err := extractBookIDFromStartQuery(inner_update.Message.Text)
				if err != nil {
					log.Printf("Error occurred during extract book id from start query for delete - %s", err.Error())
					SendUnknownError(bot_api, inner_update.FromChat().ChatConfig().ChatID)
				}
				book_title, err := db_action.GetBookTitleByID(book_id)
				if err != nil {
					log.Printf("Error occurred during get book title in delete operation - %s", err.Error())
					SendUnknownError(bot_api, inner_update.FromChat().ChatConfig().ChatID)
				}
				err = db_action.DeleteBookByID(book_id)
				if err != nil {
					log.Printf("Error occurred during delete book from database - %s", err.Error())
					SendUnknownError(bot_api, inner_update.FromChat().ChatConfig().ChatID)
				}
				msg := tgbotapi.NewMessage(inner_update.FromChat().ChatConfig().ChatID, fmt.Sprintf(BOOK_DELETED_MESSAGE, book_title))
				if _, err = bot_api.Send(msg); err != nil {
					log.Printf("Error occurred during send book successfuly deleted message")
					SendError(bot_api, inner_update.FromChat().ChatConfig().ChatID, fmt.Sprintf(BOOK_DELETED_MESSAGE, book_title))
				}
				break
				// User canceled the operation
			} else if inner_update.CallbackQuery != nil && inner_update.CallbackQuery.Data != "" && inner_update.CallbackQuery.Data == CANCEL_KEYBOARD_ITEM_TITLE {
				break
			}
		}
	}
}
func BackToMainMenu(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Call manualy start command handler
	Start_CommandHandler(bot_api, update)
}
func sendNoOrdersInConfirmationQueue(bot_api *tgbotapi.BotAPI, chat_id int) error {
	msg := tgbotapi.NewMessage(int64(chat_id), NO_ORDERS_IN_CONFIRMATION_QUEUE)
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send no orders in confirmation queue - %s", err.Error())
		return err
	}
	return nil
}
func Admin_ConfirmOrders_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) {
	// Get in confirmation queue orders
	orders, err := db_action.GetInConfirmationQueueOrders()
	if err != nil {
		log.Printf("Error occurred during get in confirmation queue orders - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
	if len(orders) < 1 {
		err = sendNoOrdersInConfirmationQueue(bot_api, int(update.FromChat().ChatConfig().ChatID))
		if err != nil {
			log.Print(err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
	}
	for i := range orders {
		// Get user telegram id
		user_telegram_id, err := db_action.GetUserTelegramIDByUserID(int(orders[i].UserID))
		if err != nil {
			log.Printf("Error occurred during get user telegram id by user id - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
		// Create msg message
		message, err := makeCartMessage(user_telegram_id)
		if err != nil {
			log.Printf("Error occurred during make cart message for confirm or reject - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
		// Create msg
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, message)
		msg.ParseMode = "html"
		msg.ReplyMarkup = CONFIRM_OR_REJECT_ORDER_INLINE_KEYBOARD
		// Send msg to user
		if _, err := bot_api.Send(msg); err != nil {
			log.Printf("Error occureed during send message to confirm or reject order - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
		// Wait for user choose
		for inner_update := range *updates {
			if inner_update.CallbackQuery != nil && inner_update.CallbackQuery.Data != "" {
				// User canceled confirmation orders
				if inner_update.CallbackQuery.Data == CANCEL_KEYBOARD_ITEM_TITLE {
					msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, CONFIRMATION_OR_REJECTION_ORDRES_CANCELED_MESSAGE)
					if _, err = bot_api.Send(msg); err != nil {
						log.Printf("Error occurred during send confirmation or rejection of orders canceled - %s", err.Error())
						SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
					}
					return
					// User confirmed order
				} else if inner_update.CallbackQuery.Data == CONFIRM_ORDER_KEYBOARD_ITEM {
					// Change order status
					err = db_action.ChangeOrderStatus(orders[i].ID, db_action.IN_PACKING_QUEUE_ORDER_STATUS_ID)
					if err != nil {
						log.Printf("Error occurred during change order status to in packing queue - %s", err.Error())
						SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
						break
					}
					// Send status of order to user
					err = sendOrderConfirmedToUser(bot_api, user_telegram_id, orders[i].ID)
					if err != nil {
						log.Printf("Error occurred during send order confirmed status to user - %s", err.Error())
						SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
						break
					}
					// Send message to admin
					msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, ORDER_CONFIRMED_MESSAGE)
					if _, err = bot_api.Send(msg); err != nil {
						log.Printf("Error occurred during send order confirmed message - %s", err.Error())
						SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
						break
					}
					// User rejected order
				} else if inner_update.CallbackQuery.Data == REJECT_ORDER_KEYBOARD_ITEM {
					err = db_action.ChangeOrderStatus(orders[i].ID, db_action.REJECTED_ORDER_STATUS_ID)
					if err != nil {
						log.Printf("Error occurred during change order status to rejected - %s", err.Error())
						SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
						break
					}
					// Send status of order to user
					err = sendOrderRejectedToUser(bot_api, user_telegram_id, orders[i].ID)
					if err != nil {
						log.Printf("Error occurred during send order rejected status to user - %s", err.Error())
						SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
						break
					}
					// Send message to admin
					msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, ORDER_REJECTED_MESSAGE)
					if _, err = bot_api.Send(msg); err != nil {
						log.Printf("Error occurred during send order rejected message - %s", err.Error())
						SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
						break
					}
				}
			}
			break
		}

	}
	// Orders in confirmation queue finished
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, THERE_IS_NO_IN_CONFIRMATION_ORDER_LEFT)
	if _, err = bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send in confirmation order queue finished message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func Admin_Statistics_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}
func Admin_BackToAdminPanel_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	ADMIN_WANTS_TO_GO_USER_MODE = false
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, ADMIN_START_TEXT)
	keyboard, err := makeMainKeyboard(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during make main keyboard - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
	msg.ReplyMarkup = keyboard
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send start message to normal user - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func Admin_BackToUserPanel_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	ADMIN_WANTS_TO_GO_USER_MODE = true
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, START_TEXT)
	keyboard, err := makeMainKeyboard(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during make main keyboard - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
	msg.ReplyMarkup = keyboard
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send start message to normal user - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func ShowUserOrders_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Make message of msg
	message, err := makeShowUserOrdersMessage(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred druing make show user orders message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, message)
	msg.ParseMode = "html"
	if _, err = bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send show user orders message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}

func FAQ_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, FAQ_MESSAGE)
	keyboard, err := makeMainKeyboard(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during make main keyboard - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
	msg.ReplyMarkup = keyboard
	msg.ParseMode = "html"
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send start message to normal user - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
