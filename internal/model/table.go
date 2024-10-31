package model

import "order_manager/internal/id"

type Table struct {
	ID       id.ID
	OrderIDs []id.ID
	BillID   id.ID
}

type Table2 struct {
	ID     id.ID
	Orders []Order
	Bill   []Bill
}
