package db_action

type CartInformationForCalculateShipmentCost struct {
	SendProvince string `gorm:"column:Province"`
	SendCity     string `gorm:"column:City"`
	BooksInfo    []BookPriceAndWeight
}
type BookPriceAndWeight struct {
	Price  float32
	Weight float32
}
