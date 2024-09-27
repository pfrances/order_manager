package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type Preparation struct {
	iD      id.ID
	orderID id.ID
	status  model.PreparationStatus
}

func newPreparation(p *model.Preparation) Preparation {
	return Preparation{
		iD:      p.ID,
		orderID: p.OrderID,
		status:  model.PreparationStatusPending,
	}
}

func (p Preparation) toModel() model.Preparation {
	return model.Preparation{
		ID:      p.iD,
		OrderID: p.orderID,
		Status:  p.status,
	}
}
