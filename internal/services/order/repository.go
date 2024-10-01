package order

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type Repository interface {
	CreateOrder(order *model.Order) error
	GetOrder(id id.ID) *model.Order
	RemoveOrder(id id.ID) error
	UpdateOrder(id id.ID, fn func(order *model.Order) error) error
}
