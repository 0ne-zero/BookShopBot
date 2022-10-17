package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var SEARCH_BOOK_INLINE_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.InlineKeyboardButton{Text: "وارد کردن عبارت حستجو", SwitchInlineQueryCurrentChat: &EMPTY_STRING}),
)
