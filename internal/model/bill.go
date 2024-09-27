package model

import (
	"order_manager/internal/id"
	"order_manager/internal/money"
)

type Bill struct {
	ID          id.ID
	TotalAmount money.Money
	AlreadyPaid money.Money
}
