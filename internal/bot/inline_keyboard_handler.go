package bot

import (
	"log"
	"strconv"
	"strings"

	db_action "github.com/0ne-zero/BookShopBot/internal/database/action"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AddBookToCart_InlineKeyboardHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	data := update.CallbackQuery.Data
	// Extract book id from  callback data (callback data = ADD_BOOK_TO_CART + "?" + <BOOK_ID>)
	book_id_str := data[strings.LastIndex(data, "?"):]
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
	err = db_action.AddBookToCart(cart_id, book_id)
	if err != nil {
		log.Printf("Error occurred during add book to cart - %s", err.Error())
		SendError(bot_api, update.FromChat().ChatConfig().ChatID, BOOK_NOT_ADDED_TO_CART_SUCCESSFULY)
	} else {
		msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, BOOK_ADDED_TO_CART_SUCCESSFULY)
		if _, err := bot_api.Send(msg); err != nil {
			log.Printf("Error occurred during send book added to cart successfuly message - %s", err.Error())
			SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
	}
}
