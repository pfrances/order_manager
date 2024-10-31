package model

import (
	"order_manager/internal/id"
	"order_manager/internal/money"
)

type BillStatus string

const (
	BillStatusPending BillStatus = "PENDING"
	BillStatusPaid    BillStatus = "PAID"
)

type Bill struct {
	ID          id.ID
	TableID     id.ID
	Items       []MenuItem
	Status      BillStatus
	TotalAmount money.Money
	AlreadyPaid money.Money
}
