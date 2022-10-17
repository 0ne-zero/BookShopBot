package main

import (
	"log"
	"path/filepath"

	"github.com/0ne-zero/BookShopBot/internal/bot"
	"github.com/0ne-zero/BookShopBot/internal/database"
	"github.com/0ne-zero/BookShopBot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var LOG_FILE_PATH = filepath.Join("../log", "log.txt")

func main() {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	err := database.MigrateModels(db)
	if err != nil {
		log.Fatalf("Cannot migrate models with database - %s", err.Error())
	}

	// Config logger
	f, err := utils.OpenLogFile(LOG_FILE_PATH)
	if err != nil {
		log.Fatalf("Error occurred during open log file - %s\n", err.Error())
	}
	log.SetOutput(f)
	log.SetFlags(log.Llongfile)

	log.Println("Starting ...")

	// Create bot
	// Config how to update messages
	bot_api, updates, err := bot.ConfigBot()

	if err != nil {
		log.Fatalf("Error occurred during config bot - %s\n", err.Error())
	}

	log.Printf("Authorized on account %s\n", bot_api.Self.UserName)

	// Updates(Events) Handler
	for update := range updates {
		// Is
		bot_api.Request(tgbotapi.NewCallback(update.InlineQuery.ID, "xx"))
		// Is admin
		if bot.IsAdmin(&update) {
			// User is admin
			if update.Message.Text != "" {
				// Is command
				if bot.IsCommand(update.Message.Text) {
					switch update.Message.Text {
					case "/start":
						bot.Admin_Start_CommandHandler(bot_api, &update)
					}
				}
				switch update.Message.Text {
				case bot.ADMIN_ADD_BOOK_KEYBOARD_ITEM_TITLE:
					bot.Admin_AddBook_KeyboardHandler(bot_api, &update)
				case bot.ADMIN_CONFIRM_ORDERS_KEYBOARD_ITEM_TITLE:
					bot.Admin_ConfirmOrders_KeyboardHandler(bot_api, &update)
				case bot.ADMIN_STATISTICS_KEYBOARD_ITEM_TITLE:
					bot.Admin_Statistics_KeyboardHandler(bot_api, &update)
				}
			}
			// User is normal user
		} else {
			// Is a inline query
			if update.InlineQuery != nil && update.InlineQuery.Query != "" {
				bot.SearchBookByTitle_KeyboardHandler(bot_api, &update)
			}
			if update.Message.Text != "" {
				// Is a command
				if bot.IsCommand(update.Message.Text) {
					switch update.Message.Text {
					case "/start":
						bot.Start_CommandHandler(bot_api, &update)
					}
				}
				switch update.Message.Text {
				case bot.SEARCH_BOOK_KEYBOARD_ITEM_TITLE:
					bot.SearchBookByTitle_KeyboardHandler(bot_api, &update)
				case bot.CART_KEYBOARD_ITEM_TITLE:
					bot.Cart_KeyboardHandler(bot_api, &update)
				case bot.BUY_CART_KEYBOARD_ITEM_TITLE:
					bot.BuyCart_KeyboardHandler(bot_api, &update)
				case bot.CONTACT_ADMIN_KEYBOARD_ITEM_TITLE:
					bot.ContactAdmin_KeyboardHandler(bot_api, &update)
				}
			}
		}
	}
}
