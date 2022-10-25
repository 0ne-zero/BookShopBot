package main

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/0ne-zero/BookShopBot/internal/bot"
	"github.com/0ne-zero/BookShopBot/internal/database"
	db_action "github.com/0ne-zero/BookShopBot/internal/database/action"
	"github.com/0ne-zero/BookShopBot/internal/database/model"
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
	bot.BOT_USERNAME = bot_api.Self.UserName
	log.Printf("Authorized on account %s\n", bot_api.Self.UserName)

	// Updates(Events) Handler
	for update := range updates {
		// Add user to database (If not exists)
		if from := update.SentFrom(); from != nil {
			is_exists, err := db_action.IsUserExistsByTelegramUserID(int(from.ID))
			if err != nil {
				log.Printf("Error occurred during check user already exists in database - %s", err.Error())
			} else if !is_exists {
				err := db_action.AddUser(&model.User{UserTelegramID: int(from.ID), UserTelegramUsername: from.UserName})
				if err != nil {
					log.Printf("Error occurred during add user to database - %s", err.Error())
				}
			}
		}
		// Public operations
		// Inline keyboard (call back query) handlers
		if update.CallbackQuery != nil && update.CallbackQuery.Data != "" {
			// Alias for call back data
			data := &update.CallbackQuery.Data
			switch {
			case strings.Contains(*data, bot.ADD_BOOK_TO_CART_INLINE_KEYBOARD_ITEM_TITLE):
				bot.AddBookToCart_InlineKeyboardHandler(bot_api, &update)
			case strings.Contains(*data, bot.DELETE_BOOK_FROM_CART_INLINE_KEYBOARD_ITEM_TITLE):
				bot.DeleteBookFromCart_InlineKeyboardHandler(bot_api, &update)
			case *data == bot.BUY_CART_KEYBOARD_ITEM_TITLE:
				bot.BuyCart_InlineKeyboardHandler(bot_api, &update)
			case *data == bot.ADDRESS_KEYBOARD_ITEM_TITLE:
				bot.Address_InlineKeyboardHandler(bot_api, &update, &updates)
			case *data == bot.CLICK_FOR_EDIT_ADDRESS_INLINE_KEYBOARD_ITEM_TITLE:
				bot.SetAddress_InlineKeyboardHandler(bot_api, &update, &updates)
			case *data == bot.CANCEL_KEYBOARD_ITEM_TITLE:
				bot.BackToMainMenu(bot_api, &update)
			case *data == bot.I_PAID_CART_INLINE_KEYBOARD_ITEM_TITLE:
				bot.IPaidCart_InlineKeyboardHandler(bot_api, &update)
			}
		}
		// Is admin
		if bot.IsAdmin(int(update.SentFrom().ID)) {
			// User is admin
			if update.InlineQuery != nil && update.InlineQuery.Query != "" {
				// Search for delete
				if strings.Contains(update.InlineQuery.Query, bot.DELETE_STRING) {
					bot.SearchBookByTitleForDelete_InlineQueryHandler(bot_api, &update)
				} else {
					bot.SearchBookByTitle_InlineQueryHandler(bot_api, &update)
				}
			}
			if update.Message != nil && update.Message.Text != "" {
				admin_Message_Text_Handler(bot_api, &update, &updates)
			}

			// User is normal user
		} else {
			// Is a inline query
			if update.InlineQuery != nil && update.InlineQuery.Query != "" {
				bot.SearchBookByTitle_InlineQueryHandler(bot_api, &update)
			}

			if update.Message != nil && update.Message.Text != "" {
				user_Message_Text_Handler(bot_api, &update, &updates)
			}
		}
	}
}

func user_Message_Text_Handler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) {
	// Is a start query (book search)
	if bot.IsStartQuery(update.Message.Text) {
		bot.StartQueryHandler(bot_api, update)
	}

	// Is a command
	if bot.IsCommand(update.Message.Text) {
		switch update.Message.Text {
		// User, start command handler
		case "/start":
			bot.Start_CommandHandler(bot_api, update)
		}
	}
	switch update.Message.Text {
	// Go to search book (inline query mode)
	case bot.SEARCH_BOOK_KEYBOARD_ITEM_TITLE:
		bot.SearchBookByTitle_KeyboardHandler(bot_api, update)
		// Show cart handler
	case bot.CART_KEYBOARD_ITEM_TITLE:
		bot.Cart_KeyboardHandler(bot_api, update)
		// Buy cart handler
	case bot.BUY_CART_KEYBOARD_ITEM_TITLE:
		bot.Cart_KeyboardHandler(bot_api, update)
		// Contact to admin handler
	case bot.CONTACT_ADMIN_KEYBOARD_ITEM_TITLE:
		bot.ContactAdmin_KeyboardHandler(bot_api, update)
		// Set address handler
	case bot.ADDRESS_KEYBOARD_ITEM_TITLE:
		bot.SetAddress_KeyboardHandler(bot_api, update, updates)
		// Back to main menu (admin panel) handler
	case bot.MAIN_MENU_ITEM_TITLE:
		bot.BackToMainMenu(bot_api, update)
	}
}
func admin_Message_Text_Handler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel) {
	// Is command
	if bot.IsCommand(update.Message.Text) {
		switch update.Message.Text {
		case "/start":
			// Admin start command handler (admin panel)
			bot.Admin_Start_CommandHandler(bot_api, update)
		}
	}
	// Handle start query
	if bot.IsStartQuery(update.Message.Text) {
		bot.StartQueryHandler(bot_api, update)
	}
	switch update.Message.Text {
	// Back to main menu (admin panel) handler
	case bot.MAIN_MENU_ITEM_TITLE:
		bot.BackToMainMenu(bot_api, update)
		// Add book handler
	case bot.ADMIN_ADD_BOOK_KEYBOARD_ITEM_TITLE:
		bot.Admin_AddBook_KeyboardHandler(bot_api, update, updates)
		// Delete book handler
	case bot.ADMIN_DELETE_BOOK_KEYBOARD_ITEM_TITLE:
		bot.Admin_DeleteBook_KeyboardHandler(bot_api, update, updates)
		// Confirm orders handler
	case bot.ADMIN_CONFIRM_ORDERS_KEYBOARD_ITEM_TITLE:
		bot.Admin_ConfirmOrders_KeyboardHandler(bot_api, update, updates)
		// Statistics handler
	case bot.ADMIN_STATISTICS_KEYBOARD_ITEM_TITLE:
		bot.Admin_Statistics_KeyboardHandler(bot_api, update)
		// Back to admin user panel handler
	case bot.ADMIN_BACK_TO_USER_PANEL_ITEM_TITLE:
		bot.Admin_BackToUserPanel_KeyboardHandler(bot_api, update)
		// Back to admin panel handler
	case bot.ADMIN_BACK_TO_ADMIN_PANEL_ITEM_TITLE:
		bot.Admin_BackToAdminPanel_KeyboardHandler(bot_api, update)
		// Go to search book (inline query mode)
	case bot.SEARCH_BOOK_KEYBOARD_ITEM_TITLE:
		bot.SearchBookByTitle_KeyboardHandler(bot_api, update)
		// Go to set address
	case bot.ADDRESS_KEYBOARD_ITEM_TITLE:
		bot.SetAddress_KeyboardHandler(bot_api, update, updates)
		// Contact to admin handler
	case bot.CONTACT_ADMIN_KEYBOARD_ITEM_TITLE:
		bot.ContactAdmin_KeyboardHandler(bot_api, update)
		// Show user orders
	case bot.SHOW_ORDERS_KEYBOARD_ITEM_TITLE:
		bot.ShowUserOrders_KeyboardHandler(bot_api, update)
	}
}
