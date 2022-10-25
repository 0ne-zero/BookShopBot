package bot

import (
	"fmt"
	"log"

	db_action "github.com/0ne-zero/BookShopBot/internal/database/action"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SearchBookByTitle_InlineQueryHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	query := update.InlineQuery.Query
	// Query is too short
	if len(query) < 2 {
		item := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, ENTERED_PHRASE_IS_TOO_SHORT_ERROR, AT_LEAST_ENTER_ONE_CHARACTER_ERROR)
		item.Description = AT_LEAST_ENTER_ONE_CHARACTER_ERROR
		result_config := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    true,
			Results:       []interface{}{item},
		}
		if _, err := bot_api.Request(result_config); err != nil {
			log.Printf("Error occurred during send query is too short for search - %s", err.Error())
		}
	} else {
		// Get books
		books, err := db_action.SearchBooksByTitle(query)
		if err != nil {
			log.Printf("Error occurred during search book by title - %s", err.Error())
		}
		books_len := len(books)
		// No result found
		if books_len == 0 {
			item := tgbotapi.NewInlineQueryResultArticle(update.InlineQuery.ID, NO_RESULT_FOUND_ERROR, fmt.Sprintf(NO_RESULT_FOUND_DESCRIPTION_FORMAT_ERROR, query))
			item.Description = fmt.Sprintf(NO_RESULT_FOUND_DESCRIPTION_FORMAT_ERROR, query)
			result_config := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				Results:       []interface{}{item},
			}
			if _, err := bot_api.Request(result_config); err != nil {
				log.Printf("Error occurred during send no found result for sent query - %s", err.Error())
			}
			// Some result found
		} else {
			result_cfg := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     10,
				Results:       []interface{}{},
			}
			for i := 0; i < books_len; i++ {
				item := tgbotapi.NewInlineQueryResultArticle(fmt.Sprint(books[i].ID), fmt.Sprintf("%s (%s)", books[i].Title, books[i].Author), fmt.Sprintf(BOT_START_QUERY, bot_api.Self.UserName, books[i].ID))
				item.Description = fmt.Sprintf("نویسنده: %s\nمترجم: %s\nدسته بندی: %s", books[i].Author, books[i].Translator, books[i].Genre)
				result_cfg.Results = append(result_cfg.Results, item)
			}
			if _, err := bot_api.Request(result_cfg); err != nil {
				log.Printf("Error occurred during send query result - %s", err.Error())
			}
		}
	}
}

// Incomplete
func SearchBookByTitleForDelete_InlineQueryHandler(bot_api *tgbotapi.BotAPI, update *tgbotapi.Update) {
	SearchBookByTitle_InlineQueryHandler(bot_api, update)
}
