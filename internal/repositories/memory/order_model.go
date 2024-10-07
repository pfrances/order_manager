package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type order struct {
	iD             id.ID
	tableID        id.ID
	itemIDs        []id.ID
	status         model.OrderStatus
	preparationIDs []id.ID
}

func orderFromModel(m model.Order) order {
	return order{
		iD:             m.ID,
		tableID:        m.TableID,
		itemIDs:        m.MenuItemIDs,
		status:         m.Status,
		preparationIDs: m.PreparationIDs,
	}
}

func (m order) toModel() model.Order {
	return model.Order{
		ID:             m.iD,
		TableID:        m.tableID,
		MenuItemIDs:    m.itemIDs,
		Status:         m.status,
		PreparationIDs: m.preparationIDs,
	}
}
