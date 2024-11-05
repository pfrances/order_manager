package inmem

import (
	"context"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"sync"
)

type Bill struct {
	bills map[id.ID]domain.Bill
	mu    sync.Mutex
}

func NewBill() *Bill {
	return &Bill{
		bills: make(map[id.ID]domain.Bill),
	}
}

func (b *Bill) Save(ctx context.Context, bill domain.Bill) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.bills[bill.ID] = bill
	return nil
}

func (b *Bill) Find(ctx context.Context, id id.ID) (domain.Bill, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	bill, ok := b.bills[id]
	if !ok {
		return domain.Bill{}, domain.Errorf(domain.ENOTFOUND, "bill with id %s not found", id)
	}
	return bill, nil
}
