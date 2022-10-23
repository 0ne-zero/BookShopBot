package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Back to main keyboard
var BACK_TO_MAIN_MENU_KEYBOARD = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(MAIN_MENU_ITEM_TITLE),
	),
)

// Add book static keyboard
var CENSORED_STATUS_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(IS_NOT_CENSORED_STATUS_KEYBOARD_ITEM_TITLE, "0"),
		tgbotapi.NewInlineKeyboardButtonData(IS_CENSORED_STATUS_KEYBOARD_ITEM_TITLE, "1"),
	),
)

var CANCEL_KEYBOARD = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(CANCEL_KEYBOARD_ITEM_TITLE),
	),
)
