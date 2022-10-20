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
	TelegramUsername string
	TelegramUserID   string
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
	Province       string
	City           string `gorm:"NOT NULL;"`
	Street         string `gorm:"NOT NULL;"`
	BuildingNumber string `gorm:"NOT NULL;"`
	PostalCode     string
	Description    string
	// Address has one User
	UserID uint `gorm:"NOT NULL;"`
}
type Book struct {
	Base
	ISBN           string
	Title          string
	Author         string
	Translator     string
	PaperType      string
	Description    string
	NumberOfPages  int
	Genre          string
	Censored       bool
	Publisher      string
	PublishDate    string
	Price          float64
	GoodReadsScore float32
	ArezoScore     float32

	// Book has many BookCoverTypes
	CoverType       *BookCoverType `gorm:"foreignkey:BookCoverTypeID"`
	BookCoverTypeID int
	// Book has many BookSize
	BookSize   *BookSize `gorm:"foreignkey:BookSizeID"`
	BookSizeID int
	// Book has one BookAgeCategory
	BookAgeCategoryID int
	BookAgeCategory   *BookAgeCategory `gorm:"foreignkey:BookAgeCategoryID"`
}
type BookCoverType struct {
	Base
	Type string
	// BookCoverType has many Book
	Books []*Book
}
type BookSize struct {
	Base
	Name string
	// BookSize has many Book
	Books []*Book
}
type BookAgeCategory struct {
	Base
	Category string
	// BookAgeCategory has many Book
	Books []*Book
}

// Ordering
type Order struct {
	Base
	// Order has one Cart
	Cart   *Cart
	CartID uint
	// Order has one User
	UserID uint
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
