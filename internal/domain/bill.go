package domain

import (
	"context"
	"order_manager/internal/id"
)

type BillStatus string

const (
	BillStatusPending BillStatus = "PENDING"
	BillPartiallyPaid BillStatus = "PARTIALLY_PAID"
	BillStatusPaid    BillStatus = "PAID"
)

type Bill struct {
	ID          id.ID
	TableID     id.ID
	Items       []MenuItem
	Status      BillStatus
	TotalAmount int
	AlreadyPaid int
}

type BillRepository interface {
	Save(ctx context.Context, bill Bill) error
	Find(ctx context.Context, id id.ID) (Bill, error)
}

type BillService struct {
	repo BillRepository
}

func NewBillService(repo BillRepository) *BillService {
	return &BillService{repo: repo}
}

func (s *BillService) GenerateBill(ctx context.Context, table Table) (Bill, error) {
	if table.Status != TableStatusClosed {
		return Bill{}, Errorf(EINVALID, "table with id %s is not closed", table.ID)
	}

	bill := Bill{
		ID:      id.NewID(),
		TableID: table.ID,
		Status:  BillStatusPending,
	}

	for _, order := range table.Orders {
		for _, preparation := range order.Preparations {
			bill.Items = append(bill.Items, preparation.MenuItem)
			bill.TotalAmount += preparation.MenuItem.Price
		}
	}

	return bill, s.repo.Save(ctx, bill)
}

func (s *BillService) PayBill(ctx context.Context, billID id.ID, amount int) error {
	bill, err := s.repo.Find(ctx, billID)
	if err != nil {
		return err
	}

	if bill.Status == BillStatusPaid {
		return Errorf(EINVALID, "bill with id %s is already paid", billID)
	}

	if bill.AlreadyPaid+amount > bill.TotalAmount {
		return Errorf(EINVALID, "amount paid is more than total amount")
	}

	bill.AlreadyPaid += amount
	if bill.AlreadyPaid == bill.TotalAmount {
		bill.Status = BillStatusPaid
	} else {
		bill.Status = BillPartiallyPaid
	}

	return s.repo.Save(ctx, bill)
}
