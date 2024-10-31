package domain_test

import (
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"order_manager/internal/inmem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBillSuccess(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)
	table := domain.Table{ID: id.NewID(), Orders: []domain.Order{
		{ID: id.NewID(), Preparations: []domain.Preparation{
			{MenuItem: domain.MenuItem{ID: id.NewID(), Name: "Spaghetti", Price: 100}},
			{MenuItem: domain.MenuItem{ID: id.NewID(), Name: "Pizza", Price: 150}},
		}},
	}}

	bill, err := billService.GenerateBill(table)

	require.NoError(t, err, "failed to generate bill")
	assert.Equal(t, domain.BillStatusPending, bill.Status, "bill status not pending")
	assert.Equal(t, 250, bill.TotalAmount, "bill total amount not correct")
}

func TestPayBillSuccess(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)
	bill := domain.Bill{
		ID:      id.NewID(),
		TableID: id.NewID(),
		Status:  domain.BillStatusPending,
		Items: []domain.MenuItem{
			{ID: id.NewID(), Name: "Spaghetti", Price: 100},
			{ID: id.NewID(), Name: "Pizza", Price: 150},
		},
		TotalAmount: 250,
	}
	billRepo.Save(bill)

	err := billService.PayBill(bill.ID, 250)

	require.NoError(t, err, "failed to pay bill")
	bill, err = billRepo.Find(bill.ID)
	require.NoError(t, err, "failed to find bill")
	assert.Equal(t, domain.BillStatusPaid, bill.Status, "bill status not paid")
	assert.Equal(t, 250, bill.AlreadyPaid, "bill already paid not correct")
}

func TestPayBillPartiallySuccess(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)
	bill := domain.Bill{
		ID:      id.NewID(),
		TableID: id.NewID(),
		Status:  domain.BillStatusPending,
		Items: []domain.MenuItem{
			{ID: id.NewID(), Name: "Spaghetti", Price: 100},
			{ID: id.NewID(), Name: "Pizza", Price: 150},
		},
		TotalAmount: 250,
	}
	billRepo.Save(bill)

	err := billService.PayBill(bill.ID, 100)

	require.NoError(t, err, "failed to pay bill")
	bill, err = billRepo.Find(bill.ID)
	require.NoError(t, err, "failed to find bill")
	assert.Equal(t, domain.BillPartiallyPaid, bill.Status, "bill status not partially paid")
	assert.Equal(t, 100, bill.AlreadyPaid, "bill already paid not correct")
}

func TestPayBillAlreadyPaid(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)
	bill := domain.Bill{
		ID:          id.NewID(),
		TableID:     id.NewID(),
		Status:      domain.BillStatusPaid,
		TotalAmount: 250,
		AlreadyPaid: 250,
	}
	billRepo.Save(bill)

	err := billService.PayBill(bill.ID, 100)

	require.Error(t, err, "paying already paid bill should fail")
}

func TestPayBillMoreThanTotalAmount(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)
	bill := domain.Bill{
		ID:          id.NewID(),
		TableID:     id.NewID(),
		Status:      domain.BillStatusPending,
		TotalAmount: 250,
	}
	billRepo.Save(bill)

	err := billService.PayBill(bill.ID, 300)

	require.Error(t, err, "paying more than total amount should fail")
}

func TestPayBillNotFound(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)

	err := billService.PayBill(id.NewID(), 100)

	require.Error(t, err, "paying not found bill should fail")
}
