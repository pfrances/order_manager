package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type order struct {
	ID      id.ID
	ItemIds []id.ID
	Status  model.OrderStatus
}

func newOrder(m model.Order) order {
	return order{
		ID:      m.ID,
		ItemIds: m.ItemIds,
		Status:  m.Status,
	}
}

func (m order) toModel() model.Order {
	return model.Order{
		ID:      m.ID,
		ItemIds: m.ItemIds,
		Status:  m.Status,
	}
}
