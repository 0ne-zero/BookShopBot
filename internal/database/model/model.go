package model

import (
	"time"
)

type Model interface {
	User | Address | Book | BookCoverType | BookSize | BookAgeCategory | Order | OrderStatus | Cart | CartItem
}
type Base struct {
	ID        int `gorm:"primary"`
	CreatedAt *time.Time
}
type User struct {
	Base
	TelegramUsername string `gorm:"NOT NULL;"`
	TelegramUserID   string `gorm:"NOT NULL;"`
	// User has one Address
	Address Address
	// User has one Cart
	Cart Cart
	// User has many Orders
	Orders []*Order
}
type Address struct {
	Base
	Country        string `gorm:"NOT NULL;"`
	Province       string `gorm:"NOT NULL;"`
	City           string `gorm:"NOT NULL;"`
	Street         string `gorm:"NOT NULL;"`
	BuildingNumber string `gorm:"NOT NULL;"`
	PostalCode     string `gorm:"NOT NULL;"`
	PhoneNumber    string `gorm:"NOT NULL;"`
	Description    string `gorm:"NOT NULL;"`
	// Address has one User
	UserID uint `gorm:"NOT NULL;"`
}
type Book struct {
	Base
	ISBN          string
	Title         string `gorm:"NOT NULL;"`
	Author        string `gorm:"NOT NULL;"`
	Translator    string `gorm:"NOT NULL;"`
	PaperType     string `gorm:"NOT NULL;"`
	Description   string `gorm:"NOT NULL;"`
	NumberOfPages int    `gorm:"NOT NULL;"`
	Genre         string `gorm:"NOT NULL;"`
	// Pictures seperated with | character
	Pictures       string
	Censored       bool    `gorm:"NOT NULL;"`
	Publisher      string  `gorm:"NOT NULL;"`
	PublishDate    string  `gorm:"NOT NULL;"`
	Price          float64 `gorm:"NOT NULL;"`
	GoodReadsScore float32 `gorm:"NOT NULL;"`
	ArezoScore     float32 `gorm:"NOT NULL;"`
	Weight         float32 `gorm:"NOT NULL;"`

	// Book has many BookCoverTypes
	CoverType       *BookCoverType `gorm:"foreignkey:BookCoverTypeID"`
	BookCoverTypeID int            `gorm:"NOT NULL;"`
	// Book has many BookSize
	BookSize   *BookSize `gorm:"foreignkey:BookSizeID"`
	BookSizeID int       `gorm:"NOT NULL;"`
	// Book has one BookAgeCategory
	BookAgeCategory   *BookAgeCategory `gorm:"foreignkey:BookAgeCategoryID"`
	BookAgeCategoryID int              `gorm:"NOT NULL;"`
}
type BookCoverType struct {
	Base
	Type string `gorm:"NOT NULL;"`
	// BookCoverType has many Book
	Books []*Book
}
type BookSize struct {
	Base
	Name string `gorm:"NOT NULL;"`
	// BookSize has many Book
	Books []*Book
}
type BookAgeCategory struct {
	Base
	Category string `gorm:"NOT NULL;"`
	// BookAgeCategory has many Book
	Books []*Book
}

// Ordering
type Order struct {
	Base
	// Order has one Cart
	Cart   *Cart
	CartID uint `gorm:"NOT NULL;"`
	// Order has one User
	UserID uint `gorm:"NOT NULL;"`
	// Order has one OrderStatus
	OrderStatusID uint `gorm:"NOT NULL;"`
}
type OrderStatus struct {
	Base
	Status string `gorm:"NOT NULL;"`
	// OrderStatus has many Order
	Orders []*Order
}
type Cart struct {
	Base
	TotalPrice float64 `gorm:"NOT NULL;"`
	IsOrdered  bool    `gorm:"NOT NULL;"`

	// Cart has one User
	UserID uint `gorm:"NOT NULL;"`

	// Cart has many CartItem
	CartItems []*CartItem `gorm:"NOT NULL;"`
}
type CartItem struct {
	Base
	// CartItem has one Product
	BookID       uint  `gorm:"NOT NULL;"`
	Book         *Book `gorm:"NOT NULL;"`
	BookQuantity uint  `gorm:"NOT NULL"`
	// CartItem has one Cart
	CartID uint
}
