package bot

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"strings"

	db_action "github.com/0ne-zero/BookShopBot/internal/database/action"
	"github.com/0ne-zero/BookShopBot/internal/utils"
	setting "github.com/0ne-zero/BookShopBot/internal/utils/settings"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type validateUserinputFunc func(*tgbotapi.Update) error

func ConfigBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel, error) {
	bot, err := tgbotapi.NewBotAPI(API_KEY)
	if err != nil {
		return nil, nil, fmt.Errorf("error occurred during create new bot instance - %w", err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 10
	updates := bot.GetUpdatesChan(u)

	cfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "/start", Description: "start the bot"},
		tgbotapi.BotCommand{Command: "/test", Description: "test the bot"},
	)
	bot.Request(cfg)
	return bot, updates, nil
}

func SendUnknownError(bot *tgbotapi.BotAPI, chat_id int64) {
	_, err := bot.Send(tgbotapi.NewMessage(chat_id, UNKNOWN_ERROR))
	if err != nil {
		log.Printf("Error occurred during send error message - %s\n", err.Error())
	}
}
func SendError(bot *tgbotapi.BotAPI, chat_id int64, error_text string) {
	_, err := bot.Send(tgbotapi.NewMessage(chat_id, error_text))
	if err != nil {
		log.Printf("Error occurred during send error message - %s\n", err.Error())
	}
}
func IsCommand(text string) bool {
	return strings.HasPrefix(text, "/")
}
func MakeBookCoverTypesKeyboard() (*tgbotapi.InlineKeyboardMarkup, error) {
	types, err := db_action.GetBookCoverTypes()
	if err != nil {
		return nil, err
	}
	var keyboard_items []tgbotapi.InlineKeyboardButton
	for i := range types {
		item := tgbotapi.NewInlineKeyboardButtonData(types[i].Type, fmt.Sprint(types[i].ID))
		keyboard_items = append(keyboard_items, item)
	}
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(keyboard_items...),
	)
	return &keyboard, nil
}
func MakeBookSizeKeyboard() (*tgbotapi.InlineKeyboardMarkup, error) {
	sizes, err := db_action.GetBookSize()
	if err != nil {
		return nil, err
	}
	var keyboard_items []tgbotapi.InlineKeyboardButton
	for i := range sizes {
		item := tgbotapi.NewInlineKeyboardButtonData(sizes[i].Name, fmt.Sprint(sizes[i].ID))
		keyboard_items = append(keyboard_items, item)
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	keyboard_items_len := len(keyboard_items)
	if keyboard_items_len < 4 {
		rows = append(rows, keyboard_items)
	} else {
		start := 0
		end := 3
		number_of_rows_divide := float64(keyboard_items_len) / float64(3)
		if utils.IsFloatNumberRound(number_of_rows_divide) {
			for i := 0; i < int(number_of_rows_divide); i++ {
				rows = append(rows, keyboard_items[start:end])
				start += 3
				end += 3
			}
		} else {
			number_of_rows_divide++
			for i := 0; i < int(number_of_rows_divide); i++ {
				if end > keyboard_items_len {
					rows = append(rows, keyboard_items[start:])
				}
				rows = append(rows, keyboard_items[start:end])
				start += 3
				end += 3
			}
		}
	}
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard, nil
}
func MakeBookAgeCategoryKeyboard() (*tgbotapi.InlineKeyboardMarkup, error) {
	categories, err := db_action.GetBookAgeCategories()
	if err != nil {
		return nil, err
	}
	var keyboard_items []tgbotapi.InlineKeyboardButton
	for i := range categories {
		item := tgbotapi.NewInlineKeyboardButtonData(categories[i].Category, fmt.Sprint(categories[i].ID))
		keyboard_items = append(keyboard_items, item)
	}
	var rows [][]tgbotapi.InlineKeyboardButton
	keyboard_items_len := len(keyboard_items)
	if keyboard_items_len < 4 {
		rows = append(rows, keyboard_items)
	} else {
		start := 0
		end := 3
		number_of_rows_divide := float64(keyboard_items_len) / float64(3)
		if utils.IsFloatNumberRound(number_of_rows_divide) {
			for i := 0; i < int(number_of_rows_divide); i++ {
				rows = append(rows, keyboard_items[start:end])
				start += 3
				end += 3
			}
		} else {
			number_of_rows_divide++
			for i := 0; i < int(number_of_rows_divide); i++ {
				if end > keyboard_items_len {
					rows = append(rows, keyboard_items[start:])
				}
				rows = append(rows, keyboard_items[start:end])
				start += 3
				end += 3
			}
		}
	}
	var keyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard, nil
}
func GetInputFromUser(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update, updates *tgbotapi.UpdatesChannel, input_request_text string, validate_func validateUserinputFunc) (string, error) {
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, input_request_text)
	msg.ReplyMarkup = MAIN_MENU_KEYBOARD
	_, err := bot_api.Send(msg)
	if err != nil {
		log.Print("Error occurred during send request for book title")
		return "", err
	}
	// Wait for user input
	for inner_update := range *updates {
		if err = validate_func(&inner_update); err != nil {
			return "", err
		}
		return inner_update.Message.Text, nil
	}
	return "", fmt.Errorf("nothing happened")
}

func IsStartQuery(text string) bool {
	if strings.Contains(text, "https://t.me/Xbookshopbot/?start=") {
		return true
	} else {
		return false
	}
}
func StartQueryHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	// Extract book id from query
	book_id, err := extractBookIDFromStartQuery(update.Message.Text)
	if err != nil {
		log.Printf("Error occurred during extract book id from start query - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
		return
	}
	// Format book information as text to send
	book_formatted_info, err := formatBookInformation(book_id)
	if err != nil {
		log.Printf("Error occurred during format book information - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
	// Send result
	msg := tgbotapi.NewMessage(update.FromChat().ChatConfig().ChatID, book_formatted_info)
	msg.ReplyMarkup = nil
	if update.Message != nil && update.Message.MessageID != 0 {
		msg.ReplyToMessageID = update.Message.MessageID
	}
	if _, err := bot_api.Send(msg); err != nil {
		log.Printf("Error occurred during send book information message - %s", err.Error())
		SendUnknownError(bot_api, update.FromChat().ChatConfig().ChatID)
	}
}
func extractBookIDFromStartQuery(query string) (int, error) {
	id_str := strings.Split(query, "=")[1]
	return strconv.Atoi(id_str)
}
func formatBookInformation(book_id int) (string, error) {
	book, err := db_action.GetBookByID(book_id)
	if err != nil {
		return "", err
	}
	formatted_info := fmt.Sprintf(
		BOOK_INFORMATION_FORMAT, book.Title, book.Author, book.Translator, book.NumberOfPages, book.Genre,
		book.Censored, book.CoverType.Type, book.BookSize.Name, book.BookAgeCategory.Category, book.GoodReadsScore,
		book.ArezoScore, book.Publisher, book.PublishDate, book.ISBN, book.Price)
	return formatted_info, nil
}
func GetGoodreadsScoreByISBN(isbn string) (string, error) {
	isbn = strings.ReplaceAll(isbn, "-", "")
	cmd := exec.Command("python3", "../internal/urils/isbn_to_goodreads_score.py", isbn)
	buffer := bytes.Buffer{}
	cmd.Stdout = &buffer
	output_bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output_bytes), nil
}
func IsAdmin(update *tgbotapi.Update) bool {

	// Get admin id from settings
	admin_id_str := setting.ReadFieldInSettingData("ADMIN_TELEGRAM_ID")
	// Parse id to int64
	admin_id, err := strconv.ParseInt(admin_id_str, 10, 64)
	if err != nil {
		return false
	}
	// Check ids are equal or not
	if getUserID(update) == admin_id {
		return true
	} else {
		return false
	}
}
func getUserID(update *tgbotapi.Update) int64 {
	if update.Message != nil && update.Message.From != nil {
		return update.Message.From.ID
	}
	if update.InlineQuery != nil && update.InlineQuery.From != nil {
		return update.InlineQuery.From.ID
	}
	return 0
}
