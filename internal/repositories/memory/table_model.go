package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type table struct {
	iD       id.ID
	orderIDs []id.ID
	billID   id.ID
}

func tableFromModel(t model.Table) table {
	return table{
		iD:       t.ID,
		orderIDs: t.OrderIDs,
		billID:   t.BillID,
	}
}

func (t table) toModel() model.Table {
	return model.Table{
		ID:       t.iD,
		OrderIDs: t.orderIDs,
		BillID:   t.billID,
	}
}
