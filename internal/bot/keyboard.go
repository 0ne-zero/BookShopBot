package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
