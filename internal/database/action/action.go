package db_action

import (
	"log"

	"github.com/0ne-zero/BookShopBot/internal/database"
	"github.com/0ne-zero/BookShopBot/internal/database/model"
)

type Models interface {
	model.User | model.Book | model.BookAgeCategory | model.BookCoverType | model.BookSize | model.Address | model.Order | model.OrderStatus | model.Cart | model.CartItem
}

func AddUser(u *model.User) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	return db.Create(u).Error
}
func GetUserCartIDByTelegramUserID(telegram_user_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var cart_id int
	err := db.Model(&model.Cart{}).Select("id").Where("telegram_user_id = ?", telegram_user_id).Scan(&cart_id).Error
	return cart_id, err
}
func DeleteBookFromCart(book_id, cart_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	cart_item := &model.CartItem{
		BookID:       uint(book_id),
		CartID:       uint(cart_id),
		BookQuantity: 1,
	}
	return db.Delete(cart_item).Error
}
func IsUserExistsByTelegramUserID(user_id string) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var exists bool
	err := db.Model(&model.User{}).Select("count(*) > 0").Where("telegram_user_id = ?", user_id).Scan(&exists).Error
	return exists, err
}
func IsBookExistsInCart(book_id, cart_id int) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var exists bool
	err := db.Model(&model.CartItem{}).Select("count(*) > 0").Where("cart_id = ? AND book_id = ?", cart_id, book_id).Scan(&exists).Error
	return exists, err
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
func AddAddress(addr *model.Address, telegram_user_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Set user id of address
	user_id, err := GetUserIDByTelegramUserID(telegram_user_id)
	if err != nil {
		return err
	}
	addr.UserID = uint(user_id)
	return db.Create(addr).Error
}
func GetUserAddressByTelegramUserID(telegram_user_id int) (*model.Address, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	user_id, err := GetUserIDByTelegramUserID(telegram_user_id)
	if err != nil {
		return nil, err
	}
	var addr *model.Address
	err = db.Model(&model.Address{}).Where("user_id = ?", user_id).Find(&addr).Error
	return addr, err
}
func GetUserCartBooksID(telegram_user_id int) ([]int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Get user cart id
	cart_id, err := GetUserCartIDByTelegramUserID(telegram_user_id)
	if err != nil {
		return nil, err
	}
	var books_id []int
	err = db.Model(&model.CartItem{}).Where("cart_id = ?", cart_id).Select("book_id").Scan(&books_id).Error
	return books_id, err
}
func DoesUserHaveAddress(telegram_user_id int) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	user_id, err := GetUserIDByTelegramUserID(telegram_user_id)
	if err != nil {
		return false, err
	}
	var has bool
	err = db.Model(&model.Address{}).Select("count(*) > 0").Where("user_id = ?", user_id).Scan(&has).Error
	return has, err
}
func GetUserIDByTelegramUserID(telegram_user_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var user_id int
	err := db.Model(&model.User{}).Select("id").Where("telegram_user_id = ?", telegram_user_id).Scan(&user_id).Error
	return user_id, err
}
func AddBookToCart(cart_id, book_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	cart_Item := &model.CartItem{
		BookID:       uint(book_id),
		CartID:       uint(cart_id),
		BookQuantity: 1,
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
func GetBookAuthorByID(b_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var author string
	err := db.Model(&model.Book{}).Where("id = ?", b_id).Select("author").Scan(&author).Error
	return author, err
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
