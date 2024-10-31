package inmem

import (
	"order_manager/internal/domain"
	"order_manager/internal/id"
)

type Bill struct {
	bills map[id.ID]domain.Bill
}

func NewBill() *Bill {
	return &Bill{
		bills: make(map[id.ID]domain.Bill),
	}
}

func (b *Bill) Save(bill domain.Bill) error {
	b.bills[bill.ID] = bill
	return nil
}

func (b *Bill) Find(id id.ID) (domain.Bill, error) {
	bill, ok := b.bills[id]
	if !ok {
		return domain.Bill{}, domain.Errorf(domain.ENOTFOUND, "bill with id %s not found", id)
	}
	return bill, nil
}
