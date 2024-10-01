package order

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(tableId id.ID, menuItemIds []id.ID) (id.ID, error) {
	order := &model.Order{
		ID:          id.NewID(),
		TableID:     tableId,
		MenuItemIDs: menuItemIds,
		Status:      model.OrderStatusTaken,
	}

	err := s.repo.CreateOrder(order)
	if err != nil {
		return id.NilID(), err
	}

	return order.ID, nil
}
