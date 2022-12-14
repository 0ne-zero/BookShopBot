package db_action

type OrderStatus int

const (
	IN_CONFIRMATION_QUEUE_ORDER_STATUS_ID OrderStatus = iota + 1
	REJECTED_ORDER_STATUS_ID
	IN_PACKING_QUEUE_ORDER_STATUS_ID
	SENDING_ORDER_STATUS_ID
	DELIVERED_ORDER_STATUS_ID
)
