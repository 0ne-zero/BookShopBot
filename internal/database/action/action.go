package db_action

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/0ne-zero/BookShopBot/internal/database"
	"github.com/0ne-zero/BookShopBot/internal/database/model"
	"github.com/0ne-zero/BookShopBot/internal/utils"
	setting "github.com/0ne-zero/BookShopBot/internal/utils/settings"
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
	if err != nil {
		return 0, err
	}
	// If cart_id is equal to 0, means that user doesn't even have cart,so we should create one,and return it's id
	if cart_id != 0 {
		return cart_id, nil
		// User doesn't have cart
	} else {
		// Create cart
		cart := &model.Cart{
			UserID: uint(user_id),
		}
		err = db.Create(cart).Error
		return cart.ID, err
	}
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
	err = db.Model(&cart).Where("user_id = ?", user_id).Preload("CartItems").Find(&cart).Error
	if err != nil {
		return false, err
	}
	// User cart is empty
	if len(cart.CartItems) < 1 {
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
	// Get current cart total price
	var current_total_price float32
	err := db.Model(&model.Cart{}).Where("id = ?", cart_id).Select("total_price").Scan(&current_total_price).Error
	if err != nil {
		return err
	}
	// Get book price
	book_price, err := GetBookPriceByBookID(book_id)
	if err != nil {
		return err
	}
	// Update cart total price
	err = db.Model(&model.Cart{}).Where("id = ?", cart_id).Update("total_price", current_total_price-book_price).Error
	if err != nil {
		return err
	}
	// Delete from cart
	var cart_item_id int
	err = db.Model(&model.CartItem{}).Select("id").Where("cart_id = ? AND book_id = ?", cart_id, book_id).Scan(&cart_item_id).Error
	if err != nil {
		return err
	}
	return db.Delete(&model.CartItem{}, cart_item_id).Error
}
func IsUserExistsByUserTelegramID(user_telegram_id int) (bool, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var exists bool
	err := db.Model(&model.User{}).Select("count(*) > 0").Where("user_telegram_id = ?", user_telegram_id).Scan(&exists).Error
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
func ConvertOrderStatusIDToOrderStatus(order_status_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var order_status string
	err := db.Model(&model.OrderStatus{}).Where("id = ?", order_status_id).Select("status").Scan(&order_status).Error
	return order_status, err
}
func getOrderBasicInfoByOrderID(order_id int) (*ShowOrder, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Get order
	var order model.Order
	err := db.Model(&model.Order{}).Where("id = ?", order_id).Preload("Cart").Find(&order).Error
	if err != nil {
		return nil, err
	}
	// Fill show order by order information
	var show_order ShowOrder
	show_order.TotalPrice = order.Cart.TotalPrice
	status, err := ConvertOrderStatusIDToOrderStatus(int(order.OrderStatusID))
	if err != nil {
		return nil, err
	}
	show_order.Status = status
	created_at, err := getOrderCreateTime(order_id)
	if err != nil {
		return nil, err
	}
	show_order.OrderedAt = created_at
	// TODO: get tracking id
	return &show_order, nil
}
func getOrderCreateTime(order_id int) (*time.Time, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var created_at *time.Time
	err := db.Model(&model.Order{}).Where("id = ?", order_id).Select("created_at").Scan(&created_at).Error
	return created_at, err
}

// Returns nil,nil, If order doesn't have any book
func getOrderBooksInfo(order_id int) ([]OrderBook, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var order model.Order
	err := db.Model(&order).Where("id = ?", order_id).Preload("Cart.CartItems.Book").Find(&order).Error
	if err != nil {
		return nil, err
	}
	if len(order.Cart.CartItems) < 1 {
		return nil, nil
	}
	var books = make([]OrderBook, len(order.Cart.CartItems))
	for i := range order.Cart.CartItems {
		item := OrderBook{
			ID:     order.Cart.CartItems[i].BookID,
			Title:  order.Cart.CartItems[i].Book.Title,
			Author: order.Cart.CartItems[i].Book.Author,
			Price:  order.Cart.CartItems[i].Book.Price,
		}
		books[i] = item
	}
	return books, nil
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
	// Get user id
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return err
	}
	// Set user id of address
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
	err := db.Model(&model.User{}).Select("id").Where("user_telegram_id = ?", user_telegram_id).Scan(&user_id).Error
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
func AddOrder(user_telegram_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Get user id
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return err
	}
	// Get user cart id
	cart_id, err := GetUserCartIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return err
	}
	// Generate Tracking code
	// Get tracking code tracking code length
	tracking_code_length_str := setting.ReadFieldInSettingData("TRACKING_CODE_LENGTH")
	tracking_code_length, err := strconv.Atoi(tracking_code_length_str)
	if err != nil {
		log.Fatalf("Connot convert setting tracking code length to int - %s", err.Error())
	}
	track_code, err := utils.GenerateRandomHex(tracking_code_length)
	if err != nil {
		return err
	}
	// Create order
	order := &model.Order{
		UserID:        uint(user_id),
		CartID:        uint(cart_id),
		OrderStatusID: uint(IN_CONFIRMATION_QUEUE_ORDER_STATUS_ID),
		TrackingCode:  track_code,
	}
	err = db.Create(order).Error
	if err != nil {
		return err
	}
	// Update cart mode to Ordered
	return db.Model(&model.Cart{}).Where("id = ?", cart_id).Update("is_ordered", true).Error
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
func GetBookPicturesPath(book_id int) ([]string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var raw_pics_path string
	err := db.Model(&model.Book{}).Where("id = ?", book_id).Select("pictures").Scan(&raw_pics_path).Error
	if err != nil {
		return nil, err
	}
	// Seperate paths with "|" character
	pics_path := strings.Split(raw_pics_path, "|")
	return pics_path, nil
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
func getUserOrdersIDByUserID(user_id int) ([]int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var ids []int
	err := db.Model(&model.Order{}).Where("user_id = ?", user_id).Select("id").Scan(&ids).Error
	return ids, err
}
func GetUserOrdersForShowByUserTelegramID(user_telegram_id int) ([]UserOrderForShow, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Get user id
	user_id, err := GetUserIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return nil, err
	}

	// Get orders id
	order_ids, err := getUserOrdersIDByUserID(user_id)
	if err != nil {
		return nil, err
	}
	// If user doesn't have any order return nil
	order_ids_len := len(order_ids)
	if order_ids_len < 1 {
		return nil, fmt.Errorf("user doesn't have any orders")
	}

	// Get orders info and fill return data
	var orders_info = make([]UserOrderForShow, order_ids_len)
	for i := range order_ids {
		// Extract order id
		order_id := order_ids[i]
		// Get order status
		order_status, err := GetOrder_OrderStatusByOrderID(order_id)
		if err != nil {
			return nil, err
		}
		order_created_at, err := getOrderCreateTime(order_id)
		if err != nil {
			return nil, err
		}
		// Get order books
		books, err := getOrderBooksInfo(order_id)
		if err != nil {
			return nil, err
		}
		// Get order tracking code
		tracking_code, err := getOrderTrackingCode(order_id)
		if err != nil {
			return nil, err
		}
		post_tracking_code, err := getOrderPostTrackingCode(order_id)
		if err != nil {
			return nil, err
		}
		item := UserOrderForShow{
			OrderTrackingCode:     tracking_code,
			OrderPostTrackingCode: post_tracking_code,
			OrderTime:             order_created_at,
			OrderStatus:           order_status,
			Books:                 books,
		}
		orders_info[i] = item
	}
	return orders_info, nil
}
func getOrderTrackingCode(order_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var track_code string
	err := db.Model(&model.Order{}).Where("id = ?", order_id).Select("tracking_code").Scan(&track_code).Error
	return track_code, err
}
func getOrderPostTrackingCode(order_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var track_code string
	err := db.Model(&model.Order{}).Where("id = ?", order_id).Select("tpost_tracking_code").Scan(&track_code).Error
	return track_code, err
}
func GetOrderByTrackingCode(tracking_code string) (*model.Order, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var order model.Order
	err := db.Model(&order).Where("tracking_code = ?", tracking_code).Find(&order).Error
	return &order, err
}
func getOrderIDByTrackingCode(tracking_code string) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var order_id int
	err := db.Model(&model.Order{}).Where("tracking_code = ?", tracking_code).Select("id").Scan(&order_id).Error
	return order_id, err
}
func GetOrderInfoByTrackingCode(tracking_code string) (*UserOrderForShow, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	// Get order id
	order_id, err := getOrderIDByTrackingCode(tracking_code)
	if err != nil {
		return nil, err
	}
	order_status, err := GetOrder_OrderStatusByOrderID(order_id)
	if err != nil {
		return nil, err
	}
	order_created_at, err := getOrderCreateTime(order_id)
	if err != nil {
		return nil, err
	}
	order_books, err := getOrderBooksInfo(order_id)
	if err != nil {
		return nil, err
	}
	post_tracking_code, err := getOrderPostTrackingCode(order_id)
	if err != nil {
		return nil, err
	}
	return &UserOrderForShow{
		OrderTrackingCode:     tracking_code,
		OrderPostTrackingCode: post_tracking_code,
		OrderTime:             order_created_at,
		OrderStatus:           order_status,
		Books:                 order_books,
	}, nil
}
func GetUserOrderInfoByOrderID(order_id int) (*UserOrderForShow, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	order_status, err := GetOrder_OrderStatusByOrderID(order_id)
	if err != nil {
		return nil, err
	}
	order_created_at, err := getOrderCreateTime(order_id)
	if err != nil {
		return nil, err
	}
	order_books, err := getOrderBooksInfo(order_id)
	if err != nil {
		return nil, err
	}
	tracking_code, err := getOrderTrackingCode(order_id)
	if err != nil {
		return nil, err
	}
	post_tracking_code, err := getOrderPostTrackingCode(order_id)
	if err != nil {
		return nil, err
	}
	return &UserOrderForShow{
		OrderTrackingCode:     tracking_code,
		OrderPostTrackingCode: post_tracking_code,
		OrderTime:             order_created_at,
		OrderStatus:           order_status,
		Books:                 order_books,
	}, nil
}
func GetOrder_OrderStatusByOrderID(order_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	order_status_id, err := GetOrderStatusID(order_id)
	if err != nil {
		return "", err
	}
	order_status, err := getOrderStatusByOrderStatusID(order_status_id)
	return order_status, err
}

// Works on OrderStatus model not Order
func getOrderStatusByOrderStatusID(order_status_id int) (string, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var order_status string
	err := db.Model(&model.OrderStatus{}).Where("id = ?", order_status_id).Select("status").Scan(&order_status).Error
	return order_status, err
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
func GetInConfirmationQueueOrders() ([]model.Order, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var orders []model.Order
	err := db.Model(&model.Order{}).Where("order_status_id = ?", IN_CONFIRMATION_QUEUE_ORDER_STATUS_ID).Find(&orders).Error
	return orders, err
}
func GetUserTelegramIDByUserID(user_id int) (int, error) {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	var user_telegram_id int
	err := db.Model(&model.User{}).Where("id = ?", user_id).Select("user_telegram_id").Scan(&user_telegram_id).Error
	return user_telegram_id, err
}
func ChangeOrderStatus(order_id int, order_status_id OrderStatus) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	err := db.Model(&model.Order{}).Where("id = ?", order_id).Update("order_status_id", order_status_id).Error
	return err
}
func EmptyUserCart(user_telegram_id int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}

	cart_id, err := GetUserCartIDByTelegramUserID(user_telegram_id)
	if err != nil {
		return err
	}
	var cart_items_id []int
	err = db.Model(model.CartItem{}).Where("cart_id = ?", cart_id).Select("id").Scan(&cart_items_id).Error
	if err != nil {
		return err
	}
	err = deleteCartItemsByID(cart_items_id)
	if err != nil {
		return err
	}
	return nil
}
func deleteCartItemsByID(cart_items_id []int) error {
	db := database.InitializeOrGetDB()
	if db == nil {
		log.Fatal("Cannot connect to the database")
	}
	for i := 0; i < len(cart_items_id); i++ {
		err := db.Delete(&model.CartItem{Base: model.Base{ID: i}}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
