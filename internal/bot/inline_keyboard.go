package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var SEARCH_BOOK_INLINE_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.InlineKeyboardButton{Text: ENTER_SEARCH_PHRASE_TEXT, SwitchInlineQueryCurrentChat: &EMPTY_STRING}),
)
var SEARCH_BOOK_FOR_DELETE_INLINE_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(tgbotapi.InlineKeyboardButton{Text: ENTER_SEARCH_PHRASE_FOR_DELETE_BOOK_TEXT, SwitchInlineQueryCurrentChat: &DELETE_STRING}),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(CANCEL_KEYBOARD_ITEM_TITLE, CANCEL_KEYBOARD_ITEM_TITLE)),
)
var FOR_EDIT_ADDRESS_INLINE_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(CLICK_FOR_EDIT_ADDRESS_INLINE_KEYBOARD_ITEM_TITLE, CLICK_FOR_EDIT_ADDRESS_INLINE_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewInlineKeyboardButtonData(CANCEL_KEYBOARD_ITEM_TITLE, CANCEL_KEYBOARD_ITEM_TITLE)),
)
var BUY_CART_INLINE_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(BUY_CART_KEYBOARD_ITEM_TITLE, BUY_CART_KEYBOARD_ITEM_TITLE)),
)
var SET_ADDRESS_INLINE_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(ADDRESS_KEYBOARD_ITEM_TITLE, ADDRESS_KEYBOARD_ITEM_TITLE)),
)
var I_PAID_CART = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(I_PAID_CART_INLINE_KEYBOARD_ITEM_TITLE, I_PAID_CART_INLINE_KEYBOARD_ITEM_TITLE),
	),
)
var CONFIRM_OR_REJECT_ORDER_INLINE_KEYBOARD = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(REJECT_ORDER, REJECT_ORDER),
		tgbotapi.NewInlineKeyboardButtonData(CONFIRM_ORDER, REJECT_ORDER),
	),
	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(CANCEL_KEYBOARD_ITEM_TITLE, CANCEL_KEYBOARD_ITEM_TITLE)),
)
