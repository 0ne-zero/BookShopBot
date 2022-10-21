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
}
func BuyCart_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}
func ContactAdmin_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}
func SetAddress_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) {
	var addr model.Address
	var err error
	var input_fetched = false
	// Get address country
	for !input_fetched {
		addr.Country, err = GetInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_COUNTRY, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address country - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}

	// Get address province
	for !input_fetched {
		addr.Province, err = GetInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_PROVINCE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address province - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get address city
	for !input_fetched {
		addr.City, err = GetInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_CITY, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address city - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get address street
	for !input_fetched {
		addr.Street, err = GetInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_STREET, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address street - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get address building number
	for !input_fetched {
		addr.BuildingNumber, err = GetInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_BUILDING_NUMBER, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address building number - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get address postal code
	for !input_fetched {
		addr.PostalCode, err = GetInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_POSTAL_CODE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address postal code - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get address description
	for !input_fetched {
		addr.Description, err = GetInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_DESCRIPTION, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address description - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
	// Get address phone number
	for !input_fetched {
		addr.PhoneNumber, err = GetInputFromUser(bot_api, update, updates, REQUEST_ADDRESS_PHONE_NUMBER, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get address phone number - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		input_fetched = true
	}
}
func Admin_AddBook_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) {
	// Create book, It should be fill with user inputs
	book := model.Book{}
	var err error
	// It's equal to true by default, jsut for enter to below for loop
	var input_fetched bool = false
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
	// Get book's title
	input_fetched = false
	for !input_fetched {
		raw_isbn, err := GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_ISBN, updateValidateFunc)
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
		book.Title, err = GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_TITLE, updateValidateFunc)
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
		book.Author, err = GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_AUTHOR, updateValidateFunc)
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
		book.Translator, err = GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_TRANSLATOR, updateValidateFunc)
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
		book.PaperType, err = GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_PAPER_TYPE, updateValidateFunc)
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
		book.Description, err = GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_DESCRIPTION, updateValidateFunc)
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
		number_of_pages_str, err := GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_NUMBER_OF_PAGES, updateValidateFunc)
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
		book.Genre, err = GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_GENRE, updateValidateFunc)
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
			log.Print("Error occurred during send request for book title")
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
		book.Publisher, err = GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_PUBLISHER, updateValidateFunc)
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
		book.PublishDate, err = GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_PUBLISHDATE, updateValidateFunc)
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
		price_str, err := GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_PRICE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		book.Price, err = strconv.ParseFloat(price_str, 32)
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
		arezo_score_str, err := GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_AZERO_SCORE, updateValidateFunc)
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
		price_str, err := GetInputFromUser(bot_api, update, updates, REQUEST_BOOK_PRICE, updateValidateFunc)
		if err != nil {
			log.Printf("Error occurred during get book author - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			continue
		}
		book.Price, err = strconv.ParseFloat(price_str, 32)
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
		SendError(bot_api, update.FromChat().ChatConfig().ChatID, BOOK_NOT_SAVED_IN_DATABASE)
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
				msg := tgbotapi.NewMessage(inner_update.FromChat().ChatConfig().ChatID, fmt.Sprintf(BOOK_DELETED_SUCCESSFULY, book_title))
				if _, err = bot_api.Send(msg); err != nil {
					log.Printf("Error occurred during send book successfuly deleted message")
					SendError(bot_api, inner_update.FromChat().ChatConfig().ChatID, fmt.Sprintf(BOOK_DELETED_SUCCESSFULY, book_title))
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
func Admin_ConfirmOrders_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}
func Admin_Statistics_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}
func Admin_BackToAdminPanel_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, ADMIN_START_TEXT)
	msg.ReplyMarkup = ADMIN_PANEL_KEYBOARD
	if _, err := bot_api.Send(msg); err != nil {
		log.Print("Error occurred during send start message to normal user")
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func Admin_BackToUserPanel_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, ADMIN_START_TEXT)
	msg.ReplyMarkup = ADMIN_USER_PANEL_KEYBOARD
	if _, err := bot_api.Send(msg); err != nil {
		log.Print("Error occurred during send start message to normal user")
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
