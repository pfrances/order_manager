package model

import "order_manager/internal/id"

type Table struct {
	ID     id.ID
	Orders []Order
	Bill   Bill
}
