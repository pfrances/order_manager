package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type order struct {
	ID             id.ID
	TableID        id.ID
	ItemIDs        []id.ID
	Status         model.OrderStatus
	PreparationIDs []id.ID
}

func orderFromModel(m model.Order) order {
	return order{
		ID:             m.ID,
		TableID:        m.TableID,
		ItemIDs:        m.MenuItemIDs,
		Status:         m.Status,
		PreparationIDs: m.PreparationIDs,
	}
}

func (m order) toModel() model.Order {
	return model.Order{
		ID:             m.ID,
		TableID:        m.TableID,
		MenuItemIDs:    m.ItemIDs,
		Status:         m.Status,
		PreparationIDs: m.PreparationIDs,
	}
}
