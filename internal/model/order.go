package model

import "order_manager/internal/id"

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "pending"
	OrderStatusReady   OrderStatus = "ready"
	OrderStatusDone    OrderStatus = "done"
)

type Order struct {
	ID      id.ID
	ItemIds []id.ID
	Status  OrderStatus
}
