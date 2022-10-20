package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start_CommandHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, START_TEXT)
	msg.ReplyMarkup = USER_PANEL_KEYBOARD
	if _, err := bot_api.Send(msg); err != nil {
		log.Print("Error occurred during send start message to normal user")
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func Admin_Start_CommandHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, ADMIN_START_TEXT)
	msg.ReplyMarkup = ADMIN_PANEL_KEYBOARD
	if _, err := bot_api.Send(msg); err != nil {
		log.Print("Error occurred during send start message to normal user")
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
