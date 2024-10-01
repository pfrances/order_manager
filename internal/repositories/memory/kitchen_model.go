package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type Preparation struct {
	iD         id.ID
	orderID    id.ID
	menuItemID id.ID
	status     model.PreparationStatus
}

func fromModel(p model.Preparation) Preparation {
	return Preparation{
		iD:         p.ID,
		orderID:    p.OrderID,
		menuItemID: p.MenuItemID,
		status:     p.Status,
	}
}

func (p Preparation) toModel() model.Preparation {
	return model.Preparation{
		ID:         p.iD,
		OrderID:    p.orderID,
		MenuItemID: p.menuItemID,
		Status:     p.status,
	}
}
