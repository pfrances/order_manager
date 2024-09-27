package memory

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type OrderRepository struct {
	order map[id.ID]order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (m *OrderRepository) CreateOrder(order *model.Order) error {
	m.order[order.ID] = newOrder(*order)
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

func (m *OrderRepository) RemoveItem(id id.ID) error {
	if _, ok := m.order[id]; !ok {
		return fmt.Errorf("item not found")
	}

	delete(m.order, id)
	return nil
}
