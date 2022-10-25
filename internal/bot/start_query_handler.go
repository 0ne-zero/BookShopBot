package bot

import (
	"log"

	db_action "github.com/0ne-zero/BookShopBot/internal/database/action"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Start query handler, gets the book id and returns its information
func StartQueryHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Send book information
	sent_msg_id, err := sendBookInformation(update, bot_api)
	if err != nil {
		log.Printf("Error occurred druing send book information - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}

	// Send Keyboard message to add book to cart or remove it from cart if exists in cart
	// Get user cart id
	cart_id, err := db_action.GetUserCartIDByTelegramUserID(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during get user cart id by telegram username - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Extract book id from query
	book_id, err := extractBookIDFromStartQuery(update.Message.Text)
	if err != nil {
		log.Printf("Error occurred during extract book id from start query - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	var msg tgbotapi.MessageConfig
	// Select keyboard for book
	// If exists in user cart, keyboard should have option to remove it from the cart
	// If doesn't exists in user cart keyboard should have option to add it to the cart
	if exists, err := db_action.IsBookExistsInCart(book_id, cart_id); err != nil {
		log.Printf("Error occurred during check book exists in cart - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
		// Book exists in cart
	} else if !exists {
		msg = tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, ADD_BOOK_FROM_CART_MESSAGE)
		msg.ReplyMarkup = makeBookKeyboard(book_id)
		// Book isn't exist in cart
	} else {
		msg = tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, REMOVE_BOOK_FROM_CART_MESSAGE)
		msg.ReplyMarkup = makeBookExistsInCartKeyboard(book_id)
	}
	// Set reply
	msg.ReplyToMessageID = sent_msg_id
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send book information message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}

// Returns sent message's id
func sendBookInformation(update *tgbotapi.Update, bot_api *tgbotapi.BotAPI) (int, error) {
	// Extract book id from query
	book_id, err := extractBookIDFromStartQuery(update.Message.Text)
	if err != nil {
		log.Printf("Error occurred during extract book id from start query - %s", err.Error())
		return 0, err
	}

	// Format book information as text to send (Caption)
	book_formatted_info, err := formatBookInformation(book_id)
	if err != nil {
		log.Printf("Error occurred during format book information - %s", err.Error())
		return 0, err
	}

	// Get book pictures path
	pics_path, err := db_action.GetBookPicturesPath(book_id)
	if err != nil {
		log.Printf("Error occurred during get book pictures path - %s", err.Error())
		return 0, err
	}

	// Craete book pictures, and set caption
	var files []interface{}
	for i := range pics_path {
		if pics_path[i] == "" {
			continue
		}
		item := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath(pics_path[i]))
		if i == 0 {
			item.Caption = book_formatted_info
		}
		files = append(files, item)
	}
	// Create message
	msg := tgbotapi.NewMediaGroup(update.FromChat().ChatConfig().ChatID, files)

	// Send book information
	if res, err := bot_api.Request(msg); err != nil {
		log.Printf("Error occurred during send book information - %s", err.Error())
		return 0, err
	} else {
		sent_msg_id, err := extractMessageIDFromTelegramRawResponse(string(res.Result))
		return sent_msg_id, err
	}
}
