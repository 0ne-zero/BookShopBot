package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start_CommandHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, START_TEXT)
	var err error
	msg.ReplyMarkup, err = makeMainKeyboard(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during make main keyboard - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	if _, err = bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send start message to normal user - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func Admin_Start_CommandHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, ADMIN_START_TEXT)
	keyboard, err := makeMainKeyboard(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during make main keyboard - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	msg.ReplyMarkup = keyboard
	if _, err = bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send start message to normal user - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
