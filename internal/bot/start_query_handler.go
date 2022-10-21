package bot

import (
	"log"

	db_action "github.com/0ne-zero/BookShopBot/internal/database/action"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Start query handler, gets the book id and returns its information
func StartQueryHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Extract book id from query
	book_id, err := extractBookIDFromStartQuery(update.Message.Text)
	if err != nil {
		log.Printf("Error occurred during extract book id from start query - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Format book information as text to send
	book_formatted_info, err := formatBookInformation(book_id)
	if err != nil {
		log.Printf("Error occurred during format book information - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Send result
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, book_formatted_info)
	// Get user cart id
	cart_id, err := db_action.GetUserCartIDByTelegramUserID(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during get user cart id by telegram username - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}

	// Select keyboard for book
	// If exists in user cart, keyboard should have option to remove it from the cart
	// If not exists in user cart keyboard should have option to add it to the cart
	if exists, err := db_action.IsBookExistsInCart(book_id, cart_id); err != nil {
		log.Printf("Error occurred during check book exists in cart - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	} else if !exists {
		msg.ReplyMarkup = makeBookKeyboard(book_id)
	} else {
		msg.ReplyMarkup = makeBookExistsInCartKeyboard(book_id)
	}
	// Set reply
	if update.Message != nil && update.Message.MessageID != 0 {
		msg.ReplyToMessageID = update.Message.MessageID
	}
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send book information message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
