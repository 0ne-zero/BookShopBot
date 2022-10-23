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
func GetUserCartIDByTelegramUserID(user_telegram_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Set user id of address
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return 0, err
	}
	var cart_id int
	err = db.Model(&model.Cart{}).Select("id").Where("user_id = ?", user_id).Scan(&cart_id).Error
	return cart_id, err
}
func IsUserCartEmptyByUserTelegramID(user_telegram_id int) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Set user id of address
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return false, err
	}
	var cart model.Cart
	err = db.Model(&cart).Where("user_id = ?", user_id).Preload("CartItems").First(&cart).Error
	if err != nil {
		return false, err
	}
	// User cart is empty
	if cart.CartItems == nil {
		return true, nil
	}
	// User cart isn't empty
	return false, nil
}
func DeleteBookFromCart(book_id, cart_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var cart_item_id int
	err := db.Model(&model.CartItem{}).Select("id").Where("cart_id = ? AND book_id = ?", cart_id, book_id).Scan(&cart_item_id).Error
	if err != nil {
		return err
	}
	return db.Delete(&model.CartItem{}, cart_item_id).Error
}
func IsUserExistsByTelegramUserID(user_telegram_id int) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Set user id of address
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return false, err
	}
	var exists bool
	err = db.Model(&model.User{}).Select("count(*) > 0").Where("id = ?", user_id).Scan(&exists).Error
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
func AddAddress(addr *model.Address, user_telegram_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Set user id of address
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return err
	}
	addr.UserID = uint(user_id)
	return db.Create(addr).Error
}
func GetUserAddressByTelegramUserID(user_telegram_id int) (*model.Address, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return nil, err
	}
	var addr *model.Address
	err = db.Model(&model.Address{}).Where("user_id = ?", user_id).Find(&addr).Error
	return addr, err
}
func GetUserCartBooksID(user_telegram_id int) ([]int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Get user cart id
	cart_id, err := GetUserCartIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return nil, err
	}
	var books_id []int
	err = db.Model(&model.CartItem{}).Where("cart_id = ?", cart_id).Select("book_id").Scan(&books_id).Error
	return books_id, err
}
func DoesUserHaveAddress(user_telegram_id int) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return false, err
	}
	var has bool
	err = db.Model(&model.Address{}).Select("count(*) > 0").Where("user_id = ?", user_id).Scan(&has).Error
	return has, err
}
func GetUserIDByTelegramUserID(user_telegram_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var user_id int
	err := db.Model(&model.User{}).Select("id").Where("telegram_user_id = ?", user_telegram_id).Scan(&user_id).Error
	return user_id, err
}
func AddBookToCart(cart_id, book_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Add book to cart
	cart_Item := &model.CartItem{
		BookID:       uint(book_id),
		CartID:       uint(cart_id),
		BookQuantity: 1,
	}
	err := db.Create(cart_Item).Error
	if err != nil {
		return err
	}
	// Get current cart total price
	var current_total_price float32
	err = db.Model(&model.Cart{}).Where("id = ?", cart_id).Select("total_price").Scan(&current_total_price).Error
	if err != nil {
		return err
	}
	// Get book price
	book_price, err := GetBookPriceByBookID(book_id)
	if err != nil {
		return err
	}
	// Update cart total price
	return db.Model(&model.Cart{}).Where("id = ?", cart_id).Update("total_price", current_total_price+book_price).Error
}
func GetBookPriceByBookID(book_id int) (float32, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var price float32
	err := db.Model(&model.Book{}).Where("id = ?", book_id).Select("price").Scan(&price).Error
	return price, err
}
func GetBookByID(book_id int) (*model.Book, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var b model.Book
	err := db.Model(&model.Book{}).Preload("CoverType").Preload("BookSize").Preload("BookAgeCategory").Where("id = ?", book_id).Find(&b).Error
	return &b, err
}
func GetBookTitleByID(book_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var title string
	err := db.Model(&model.Book{}).Where("id = ?", book_id).Select("title").Scan(&title).Error
	return title, err
}
func GetUserCartTotalPrice(user_telegram_id int) (float32, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Get user cart id
	cart_id, err := GetUserCartIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return 0, err
	}
	var total_price float32
	err = db.Model(&model.Cart{}).Where("id = ?", cart_id).Select("total_price").Scan(&total_price).Error
	return total_price, err
}
func DoesUserHaveOrder(user_telegram_id int) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Get user cart id
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return false, err
	}
	var have bool
	err = db.Model(&model.Order{}).Select("count(*) > 0").Where("user_id = ?", user_id).Scan(&have).Error
	return have, err
}
func GetCartInformationForCalculateShipmentCost(user_telegram_id int) (*CartInformationForCalculateShipmentCost, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Get user cart id
	cart_id, err := GetUserCartIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return nil, err
	}
	// Get user cart books id
	var books_id []int
	err = db.Model(&model.CartItem{}).Where("cart_id = ?", cart_id).Select("book_id").Scan(&books_id).Error
	if err != nil {
		return nil, err
	}
	var books_info []BookPriceAndWeight
	// Get books price and weight
	err = db.Model(&model.Book{}).Where("id IN ?", books_id).Select("price", "weight").Scan(&books_info).Error
	if err != nil {
		return nil, err
	}
	// Get user id
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return nil, err
	}
	// Get address info
	var total_info CartInformationForCalculateShipmentCost
	err = db.Model(&model.Address{}).Where("user_id = ?", user_id).Select("province", "city").Scan(&total_info).Error
	if err != nil {
		return nil, err
	}
	total_info.BooksInfo = books_info
	return &total_info, nil
}
func GetBookAuthorByID(book_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var author string
	err := db.Model(&model.Book{}).Where("id = ?", book_id).Select("author").Scan(&author).Error
	return author, err
}
func DeleteBookByID(book_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	return db.Delete(&model.Book{}, book_id).Error
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
