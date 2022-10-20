package db_action

import (
	"log"

	"github.com/0ne-zero/BookShopBot/internal/database"
	"github.com/0ne-zero/BookShopBot/internal/database/model"
)

type Models interface {
	model.User | model.Book | model.BookAgeCategory | model.BookCoverType | model.BookSize | model.Address | model.Order | model.OrderStatus | model.Cart | model.CartItem
}

func AddBook(b *model.Book) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	return db.Create(b).Error
}
func GetOrderStatusID(order_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var status_id int
	err := db.Model(&model.Order{}).Where("id = ?", order_id).Select("order_status_id").Scan(&status_id).Error
	return status_id, err
}
func GetNotConfirmedOrders() ([]*model.Order, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var orders []*model.Order
	err := db.Model(&model.Order{}).Where("order_status_id = ?", 1).Find(&orders).Error
	return orders, err
}

func AddBookToCart(cart_id, book_id, book_qunantity int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	cart_Item := &model.CartItem{
		BookID:       uint(book_id),
		CartID:       uint(cart_id),
		BookQuantity: uint(book_qunantity),
	}
	return db.Create(cart_Item).Error
}
func GetBookByID(b_id int) (*model.Book, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var b model.Book
	err := db.Model(&model.Book{}).Preload("CoverType").Preload("BookSize").Preload("BookAgeCategory").Where("id = ?", b_id).Find(&b).Error
	return &b, err
}
func GetBookTitleByID(b_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var title string
	err := db.Model(&model.Book{}).Where("id = ?", b_id).Select("title").Scan(&title).Error
	return title, err
}
func DeleteBookByID(b_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	return db.Delete(&model.Book{}, b_id).Error
}
func SearchBooksByTitle(title string) ([]*model.Book, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var books []*model.Book
	err := db.Model(&model.Book{}).Where("title LIKE ?", title).Limit(50).Find(&books).Error
	return books, err
}
func GetBookCoverTypes() ([]*model.BookCoverType, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var types []*model.BookCoverType
	err := db.Find(&types).Error
	return types, err
}
func GetBookAgeCategories() ([]*model.BookAgeCategory, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var categories []*model.BookAgeCategory
	err := db.Find(&categories).Error
	return categories, err
}
func GetBookSize() ([]*model.BookSize, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var sizes []*model.BookSize
	err := db.Find(&sizes).Error
	return sizes, err
}

func GetBookCoverTypeByID(id int) (*model.BookCoverType, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var cover_type model.BookCoverType
	err := db.Model(&cover_type).Where("id = ?", id).Find(&cover_type).Error
	return &cover_type, err
}
func GetBookSizeByID(id int) (*model.BookSize, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var book_size model.BookSize
	err := db.Model(&book_size).Where("id = ?", id).Find(&book_size).Error
	return &book_size, err
}
func GetBookAgeCategoryByID(id int) (*model.BookAgeCategory, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var category model.BookAgeCategory
	err := db.Model(&category).Where("id = ?", id).Find(&category).Error
	return &category, err
}
