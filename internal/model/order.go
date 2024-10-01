package model

import "order_manager/internal/id"

type OrderStatus string

const (
	OrderStatusTaken   OrderStatus = "taken"
	OrderStatusDone    OrderStatus = "done"
	OrderStatusAborted OrderStatus = "aborted"
)

type Order struct {
	ID             id.ID
	TableID        id.ID
	MenuItemIDs    []id.ID
	Status         OrderStatus
	PreparationIDs []id.ID
}
