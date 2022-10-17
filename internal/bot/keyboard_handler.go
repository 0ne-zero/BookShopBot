package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SearchBookByTitle_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, SEARCH_TEXT)
	msg.ReplyMarkup = SEARCH_BOOK_INLINE_KEYBOARD
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot_api.Send(msg); err != nil {
		log.Print("Error occurred during send search book by title message")
		SendError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func Cart_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}
func BuyCart_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}
func ContactAdmin_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}

func Admin_AddBook_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}
func Admin_ConfirmOrders_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}
func Admin_Statistics_KeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
}
