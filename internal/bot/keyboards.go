package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	SEARCH_BOOK_KEYBOARD_ITEM_TITLE   = "جست و جوی بر اساس عنوان کتاب"
	CART_KEYBOARD_ITEM_TITLE          = "سبد خرید"
	BUY_CART_KEYBOARD_ITEM_TITLE      = "خرید سبد"
	CONTACT_ADMIN_KEYBOARD_ITEM_TITLE = "ارتباط با ادمین"

	ADMIN_ADD_BOOK_KEYBOARD_ITEM_TITLE       = "اضافه کردن کتاب"
	ADMIN_CONFIRM_ORDERS_KEYBOARD_ITEM_TITLE = "تأیید سفارشات"
	ADMIN_STATISTICS_KEYBOARD_ITEM_TITLE     = "امار ربات"
)

var START_KEYBOARD = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(BUY_CART_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(CART_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(SEARCH_BOOK_KEYBOARD_ITEM_TITLE),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(CONTACT_ADMIN_KEYBOARD_ITEM_TITLE),
	),
)
var ADMIN_START_KEYBOARD = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ADMIN_STATISTICS_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(ADMIN_CONFIRM_ORDERS_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(ADMIN_ADD_BOOK_KEYBOARD_ITEM_TITLE),
	),
)
