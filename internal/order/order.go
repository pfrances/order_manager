package order

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type Service struct {
	repo Repository
}

type Repository interface {
	CreateOrder(order *model.Order) error
	GetOrder(id id.ID) *model.Order
	RemoveOrder(id id.ID) error
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateOrder(ids []id.ID) model.Order {
	return model.Order{
		ID:      id.NewID(),
		ItemIds: ids,
		Status:  model.OrderStatusPending,
	}
}
