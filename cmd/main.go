package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/0ne-zero/BookShopBot/internal/bot"
	"github.com/0ne-zero/BookShopBot/internal/database"
	"github.com/0ne-zero/BookShopBot/internal/utils"
)

var LOG_FILE_PATH = filepath.Join("../log", "log.txt")
var PICTURES_DIRECTORY = filepath.Join("../pictures/books/")

func main() {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	err := database.MigrateModels(db)
	if err != nil {
		log.Fatalf("Cannot migrate models with database - %s", err.Error())
	}
	// Create essential data
	err = database.CreateEssensialData(db)
	if err != nil {
		log.Fatalf("Cannot create essential data")
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

		// Is admin
		if bot.IsAdmin(&update) {
			// User is admin
			if update.InlineQuery != nil && update.InlineQuery.Query != "" {
				// Search for delete
				if strings.Contains(update.InlineQuery.Query, bot.DELETE_STRING) {
					bot.SearchBookByTitleForDelete_InlineQueryHandler(bot_api, &update)
				} else {
					bot.SearchBookByTitle_InlineQueryHandler(bot_api, &update)
				}
			}
			if update.Message == nil {
				continue
			}
			if update.Message.Text != "" {
				// Is command
				if bot.IsCommand(update.Message.Text) {
					switch update.Message.Text {
					case "/start":
						// Admin start command handler (admin panel)
						bot.Admin_Start_CommandHandler(bot_api, &update)
					}
				}
				// Handle start query
				if bot.IsStartQuery(update.Message.Text) {
					bot.StartQueryHandler(bot_api, &update)
				}
				switch update.Message.Text {
				// Back to main menu (admin panel) handler
				case bot.MAIN_MENU_ITEM_TITLE:
					bot.BackToMainMenu(bot_api, &update)
					// Add book handler
				case bot.ADMIN_ADD_BOOK_KEYBOARD_ITEM_TITLE:
					bot.Admin_AddBook_KeyboardHandler(bot_api, &update, &updates)
					// Delete book handler
				case bot.ADMIN_DELETE_BOOK_KEYBOARD_ITEM_TITLE:
					bot.Admin_DeleteBook_KeyboardHandler(bot_api, &update, &updates)
					// Confirm orders handler
				case bot.ADMIN_CONFIRM_ORDERS_KEYBOARD_ITEM_TITLE:
					bot.Admin_ConfirmOrders_KeyboardHandler(bot_api, &update)
					// Statistics handler
				case bot.ADMIN_STATISTICS_KEYBOARD_ITEM_TITLE:
					bot.Admin_Statistics_KeyboardHandler(bot_api, &update)
					// Back to admin user panel handler
				case bot.ADMIN_BACK_TO_USER_PANEL_ITEM_TITLE:
					bot.Admin_BackToUserPanel_KeyboardHandler(bot_api, &update)
					// Back to admin panel handler
				case bot.ADMIN_BACK_TO_ADMIN_PANEL_ITEM_TITLE:
					bot.Admin_BackToAdminPanel_KeyboardHandler(bot_api, &update)
					// Go to search book (inline query mode)
				case bot.SEARCH_BOOK_KEYBOARD_ITEM_TITLE:
					bot.SearchBookByTitle_KeyboardHandler(bot_api, &update)
				}
			}
			// User is normal user
		} else {
			// Is a inline query
			if update.InlineQuery != nil && update.InlineQuery.Query != "" {
				bot.SearchBookByTitle_InlineQueryHandler(bot_api, &update)
			}
			if update.Message == nil {
				continue
			}
			if update.Message.Text != "" {
				// Is a command
				if bot.IsCommand(update.Message.Text) {
					switch update.Message.Text {
					// User start command handler
					case "/start":
						bot.Start_CommandHandler(bot_api, &update)
					}
				}
				switch update.Message.Text {
				// Go to search book (inline query mode)
				case bot.SEARCH_BOOK_KEYBOARD_ITEM_TITLE:
					bot.SearchBookByTitle_KeyboardHandler(bot_api, &update)
					// Show cart handler
				case bot.CART_KEYBOARD_ITEM_TITLE:
					bot.Cart_KeyboardHandler(bot_api, &update)
					// Buy cart handler
				case bot.BUY_CART_KEYBOARD_ITEM_TITLE:
					bot.BuyCart_KeyboardHandler(bot_api, &update)
					// Contact to admin handler
				case bot.CONTACT_ADMIN_KEYBOARD_ITEM_TITLE:
					bot.ContactAdmin_KeyboardHandler(bot_api, &update)
				}
			}
		}
	}
}
