package bot

import (
	"fmt"
	"log"
	"strconv"

	"strings"

	setting "github.com/0ne-zero/BookShopBot/internal/utils/settings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

func SendError(bot *tgbotapi.BotAPI, chat_id int64) {
	// U0001F91B = fist emoji
	err_str := "مشکلی پیش اومد, دوباره امتحان کن \U0001F91B \U0001F91B"
	_, err := bot.Send(tgbotapi.NewMessage(chat_id, err_str))
	if err != nil {
		log.Printf("Error occurred during send error message - %s\n", err.Error())
	}
}
func IsCommand(text string) bool {
	return strings.HasPrefix(text, "/")
}
func IsAdmin(update *tgbotapi.Update) bool {

	// Get admin id from settings
	admin_id_str := setting.ReadFieldInSettingData("ADMIN_TELEGRAM_ID")
	// Parse id to int64
	admin_id, err := strconv.ParseInt(admin_id_str, 10, 64)
	if err != nil {
		return false
	}
	// Check ids are equal or not
	if getUserID(update) == admin_id {
		return true
	} else {
		return false
	}
}
func getUserID(update *tgbotapi.Update) int64 {
	if update.Message != nil && update.Message.From != nil {
		return update.Message.From.ID
	}
	if update.InlineQuery != nil && update.InlineQuery.From != nil {
		return update.InlineQuery.From.ID
	}
	return 0
}
