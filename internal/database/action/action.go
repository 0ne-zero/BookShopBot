package db_action

import (
	"log"

	"github.com/0ne-zero/BookShopBot/internal/database"
	"github.com/0ne-zero/BookShopBot/internal/database/model"
)

func SearchBooksByTitle(title string) ([]*model.Book, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var books []*model.Book
	err := db.Model(&model.Book{}).Where("title LIKE ?", title).Limit(50).Find(&books).Error
	return books, err
}
