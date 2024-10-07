package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type OrderRepository struct {
	order map[id.ID]order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{order: make(map[id.ID]order)}
}

func (m *OrderRepository) CreateOrder(order model.Order) error {
	if _, ok := m.order[order.ID]; ok {
		return repositories.ErrAlreadyExists
	}

	m.order[order.ID] = orderFromModel(order)
	return nil
}

func (m *OrderRepository) GetOrder(id id.ID) *model.Order {
	order, ok := m.order[id]
	if !ok {
		return nil
	}

	modelOrder := order.toModel()
	return &modelOrder
}

func (m *OrderRepository) RemoveOrder(id id.ID) error {
	if _, ok := m.order[id]; !ok {
		return repositories.ErrNotFound
	}

	delete(m.order, id)
	return nil
}

func (m *OrderRepository) UpdateOrder(id id.ID, fn func(order *model.Order) error) error {
	order, ok := m.order[id]
	if !ok {
		return repositories.ErrNotFound
	}

	modelOrder := order.toModel()
	err := fn(&modelOrder)
	if err != nil {
		return err
	}

	m.order[id] = orderFromModel(modelOrder)
	return nil
}
