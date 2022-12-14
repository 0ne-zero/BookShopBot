package bot

import (
	"fmt"
	"log"
	"strconv"

	db_action "github.com/0ne-zero/BookShopBot/internal/database/action"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddBookToCart_InlineKeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Extract book id
	book_id_str := extractBookIDFromCallbackData(update.CallbackData())
	book_id, err := strconv.Atoi(book_id_str)
	if err != nil {
		log.Printf("Error occurred during extract book id from callback data - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Get user cart id
	cart_id, err := db_action.GetUserCartIDByTelegramUserID(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during get user cart id by their telegram username - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Add book to cart
	err = db_action.AddBookToCart(cart_id, book_id)
	if err != nil {
		log.Printf("Error occurred during add book to cart - %s", err.Error())
		SendError(bot_api, update.FromChat().ChatConfig().ChatID, BOOK_NOT_ADDED_TO_CART_MESSAGE)
		return
	}
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, BOOK_ADDED_TO_CART_MESSAGE)
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send book added to cart successfuly message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func DeleteBookFromCart_InlineKeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Extract book id from callback data
	book_id_str := extractBookIDFromCallbackData(update.CallbackData())
	book_id, err := strconv.Atoi(book_id_str)
	if err != nil {
		log.Printf("Error occurred during extract book id from callback data - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Get user cart id
	cart_id, err := db_action.GetUserCartIDByTelegramUserID(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during get user cart id by their telegram username - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Delete book from cart
	err = db_action.DeleteBookFromCart(book_id, cart_id)
	if err != nil {
		log.Printf("Error occurred during delete book from cart - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	book_name, err := db_action.GetBookTitleByID(book_id)
	if err != nil {
		log.Printf("Error occurred during get book title by id - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, fmt.Sprintf(BOOK_DELETED_FROM_CART_FORMAT, book_name))
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send book deleted form cart message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func BuyCart_InlineKeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Check does user have address
	if have_address, err := db_action.DoesUserHaveAddress(int(update.SentFrom().ID)); err != nil {
		log.Printf("Error occurred during checking does user have address - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
		// User have address
	} else if have_address {
		message, err := makeBuyCartMessage(int(update.SentFrom().ID))
		if err != nil {
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
			log.Fatalf("Error occurred during make buy cart message - %s", err.Error())
		}
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, message)
		msg.ParseMode = "html"
		msg.ReplyMarkup = I_PAID_CART
		if _, err := bot_api.Send(msg); err != nil {
			log.Printf("Error occurred during send buy cart message to user - %s", err.Error())
			SendUnknownError(bot_api, update.CallbackQuery.Message.MigrateFromChatID)
		}
		// User doesn't have address
	} else {
		// Send user that you haven't address and you should set one
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, YOU_HAVE_NOT_ADDRESS_INLINE_KEYBOARD_MESSAGE)
		msg.ReplyMarkup = SET_ADDRESS_INLINE_KEYBOARD
		if _, err := bot_api.Send(msg); err != nil {
			log.Printf("Error occurred during send you don't have address message - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
	}
}
func Address_InlineKeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) {
	SetAddress_KeyboardHandler(bot_api, update, updates)
}
func SetAddress_InlineKeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) {
	// Get user address from user
	addr, err := getUserAddressInformationFromUser(bot_api, update, updates)
	if err != nil {
		log.Printf("Error occurred during get user address form user - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Add address to user address
	err = db_action.AddAddress(addr, int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during add user address - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Send address saved message to user
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, ADDRESS_ADDED_MESSAGE)
	if _, err = bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send address addded/setted message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func IPaidCart_InlineKeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Add order
	err := db_action.AddOrder(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred druing add order - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Send message
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, fmt.Sprintf(ORDER_ADDED_FORMAT, BOT_USERNAME))
	keyboard, err := makeMainKeyboard(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during make main keyboard - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	msg.ReplyMarkup = keyboard
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send order added message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Empty user cart
	err = db_action.EmptyUserCart(int(update.SentFrom().ID))
	if err != nil {
		log.Printf("Error occurred during emtpy user cart - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
