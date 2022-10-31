package db_action

import "time"

type CartInformationForCalculateShipmentCost struct {
	SendProvince string               `gorm:"column:Province"`
	SendCity     string               `gorm:"column:City"`
	BooksInfo    []BookPriceAndWeight `gorm:"-"`
}
type BookPriceAndWeight struct {
	Price  float32
	Weight float32
}
type ShowOrder struct {
	OrderedAt    *time.Time
	TotalPrice   float32
	Status       string
	TrackingCode string
	Books        []OrderBook
}
type OrderBook struct {
	ID     uint
	Title  string
	Author string
	Price  float32
}
type UserOrderForShow struct {
	OrderTime             *time.Time
	OrderStatus           string
	OrderTrackingCode     string
	OrderPostTrackingCode string
	Books                 []OrderBook
}
type Statistics struct {
	NumberOfBooks                     uint
	NumberOfOrders                    uint
	NumberOfUsers                     uint
	NumberOfDeliveredOrders           uint
	NumberOfInConfirmationQueueOrders uint
	NumberOfRejectedOrders            uint
}
