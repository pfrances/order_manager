package bill

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type GenerateBill struct {
}

func NewGenerateBill() *GenerateBill {
	return &GenerateBill{}
}

func (g *GenerateBill) Execute(tableID id.ID) (model.Bill, error) {

	bill := model.Bill{
		ID: id.NewID(),
	}

	return bill, nil
}
