package bot

import (
	"log"

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
	}
	// Send result
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, book_formatted_info)
	msg.ReplyMarkup = nil
	if update.Message != nil && update.Message.MessageID != 0 {
		msg.ReplyToMessageID = update.Message.MessageID
	}
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send book information message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
