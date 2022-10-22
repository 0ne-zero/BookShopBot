package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	SEARCH_BOOK_KEYBOARD_ITEM_TITLE   = "جستجوی کتاب"
	CART_KEYBOARD_ITEM_TITLE          = "سبد خرید"
	BUY_CART_KEYBOARD_ITEM_TITLE      = "خرید سبد"
	CONTACT_ADMIN_KEYBOARD_ITEM_TITLE = "ارتباط با ادمین"
	SET_ADDRESS_KEYBOARD_ITEM_TITLE   = "تنظیم ادرس"
	ORDERS_KEYBOARD_ITEM_TITLE        = "سفارشات"

	MAIN_MENU_ITEM_TITLE       = "منو اصلی"
	CANCEL_KEYBOARD_ITEM_TITLE = "انصراف"

	ADMIN_ADD_BOOK_KEYBOARD_ITEM_TITLE         = "اضافه کردن کتاب"
	ADMIN_DELETE_BOOK_KEYBOARD_ITEM_TITLE      = "حدف کتاب"
	ADMIN_CONFIRM_ORDERS_KEYBOARD_ITEM_TITLE   = "تأیید سفارشات"
	ADMIN_STATISTICS_KEYBOARD_ITEM_TITLE       = "امار ربات"
	ADMIN_BACK_TO_USER_PANEL_ITEM_TITLE        = "پنل کاربر"
	ADMIN_BACK_TO_ADMIN_PANEL_ITEM_TITLE       = "پنل ادمین"
	IS_CENSORED_STATUS_KEYBOARD_ITEM_TITLE     = "سانسور شده است"
	IS_NOT_CENSORED_STATUS_KEYBOARD_ITEM_TITLE = "بدون سانسور است"
)

// User keyboards
var MAIN_MENU_KEYBOARD = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(MAIN_MENU_ITEM_TITLE),
	),
)
var USER_PANEL_KEYBOARD = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(CART_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(SET_ADDRESS_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(SEARCH_BOOK_KEYBOARD_ITEM_TITLE),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(CONTACT_ADMIN_KEYBOARD_ITEM_TITLE),
	),
)

// Admin keyboards
var ADMIN_USER_PANEL_KEYBOARD = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(CART_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(SET_ADDRESS_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(SEARCH_BOOK_KEYBOARD_ITEM_TITLE),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(CONTACT_ADMIN_KEYBOARD_ITEM_TITLE),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ADMIN_BACK_TO_ADMIN_PANEL_ITEM_TITLE),
	),
)
var ADMIN_PANEL_KEYBOARD = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ADMIN_STATISTICS_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(ADMIN_CONFIRM_ORDERS_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(ADMIN_DELETE_BOOK_KEYBOARD_ITEM_TITLE),
		tgbotapi.NewKeyboardButton(ADMIN_ADD_BOOK_KEYBOARD_ITEM_TITLE),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(ADMIN_BACK_TO_USER_PANEL_ITEM_TITLE)),
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
