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
		item := tgbotapi.InlineQueryResultArticle{
			ID:          update.InlineQuery.ID,
			Title:       ENTERED_PHRASE_IS_TOO_SHORT,
			Description: AT_LEAST_ENTER_ONE_CHARACTER,
		}
		result_config := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    true,
			Results:       []interface{}{item},
		}
		if _, err := bot_api.Request(result_config); err != nil {
			SendError(bot_api, update.FromChat().ChatConfig().ChatID)
		}
	} else {
		// Get books
		books, err := db_action.SearchBooksByTitle(query)
		if err != nil {
			log.Print("Error occurred during search book by title")
			SendError(bot_api, update.FromChat().ID)
		}
		books_len := len(books)
		// No result found
		if books_len == 0 {
			item := tgbotapi.InlineQueryResultArticle{
				ID:          update.InlineQuery.ID,
				Title:       NO_RESULT_FOUND,
				Description: fmt.Sprintf("برای عبارت %s نتیجه ای یافت نشد.\nعنوان کتاب را بررسی کنید, همچنین امکان دارد کتاب موجود نباشد.", query),
			}
			result_config := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				Results:       []interface{}{item},
			}
			if _, err := bot_api.Request(result_config); err != nil {
				SendError(bot_api, update.FromChat().ChatConfig().ChatID)
			}
			// Some result found
		} else {
			var results = make([]tgbotapi.InlineQueryResultArticle, books_len)
			for i := 0; i > books_len; i++ {
				censor_state := "بدون سانسور"
				if books[i].Censored {
					censor_state = "سانسور شده"
				}
				item := tgbotapi.InlineQueryResultArticle{
					ID:          update.InlineQuery.ID,
					Title:       fmt.Sprintf("%s (%s)", books[i].Title, books[i].PublishDate.Format("2000-01-01")),
					Description: fmt.Sprintf("نویسنده: %s\nمترجم: %s\nدسته بندی: %s\nوضعیت سانسور: %s\nانتشارات: %s", books[i].Author, books[i].Translator, books[i].Genre, censor_state, books[i].Publisher),
				}
				results = append(results, item)
			}
			result_cfg := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     10,
				Results:       []interface{}{results},
			}
			if _, err := bot_api.Request(result_cfg); err != nil {
				SendError(bot_api, update.FromChat().ChatConfig().ChatID)
			}
		}
	}
}
