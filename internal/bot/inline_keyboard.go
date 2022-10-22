package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var SEARCH_BOOK_INLINE_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.InlineKeyboardButton{Text: ENTER_SEARCH_PHRASE_TEXT, SwitchInlineQueryCurrentChat: &EMPTY_STRING}),
)
var SEARCH_BOOK_FOR_DELETE_INLINE_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.InlineKeyboardButton{Text: ENTER_SEARCH_PHRASE_FOR_DELETE_BOOK_TEXT, SwitchInlineQueryCurrentChat: &DELETE_STRING}),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(CANCEL_KEYBOARD_ITEM_TITLE, CANCEL_KEYBOARD_ITEM_TITLE)),
)
var FOR_EDIT_ADDRESS_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(CLICK_FOR_EDIT_ADDRESS_INLINE_KEYBOARD_ITEM_TITLE, CLICK_FOR_EDIT_ADDRESS_INLINE_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewInlineKeyboardButtonData(CANCEL_KEYBOARD_ITEM_TITLE, CANCEL_KEYBOARD_ITEM_TITLE),
	))
