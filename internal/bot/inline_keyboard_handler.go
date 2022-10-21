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
		SendError(bot_api, update.FromChat().ChatConfig().ChatID, BOOK_NOT_ADDED_TO_CART)
		return
	}
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, BOOK_ADDED_TO_CART_SUCCESSFULY)
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
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, fmt.Sprintf(BOOK_DELETED_FROM_CART, book_name))
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send book deleted form cart message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
