package domain_test

import (
	"context"
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
	table := domain.Table{
		ID:     id.New(),
		Status: domain.TableStatusClosed,
		Orders: []domain.Order{
			{ID: id.New(), Preparations: []domain.Preparation{
				{MenuItem: domain.MenuItem{ID: id.New(), Name: "Spaghetti", Price: 100}},
				{MenuItem: domain.MenuItem{ID: id.New(), Name: "Pizza", Price: 150}},
			}},
		},
	}

	bill, err := billService.GenerateBill(context.Background(), table)

	require.NoError(t, err, "failed to generate bill")
	assert.Equal(t, domain.BillStatusPending, bill.Status, "bill status not pending")
	assert.Equal(t, 250, bill.TotalAmount, "bill total amount not correct")
}

func TestPayBillSuccess(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)
	bill := domain.Bill{
		ID:      id.New(),
		TableID: id.New(),
		Status:  domain.BillStatusPending,
		Items: []domain.MenuItem{
			{ID: id.New(), Name: "Spaghetti", Price: 100},
			{ID: id.New(), Name: "Pizza", Price: 150},
		},
		TotalAmount: 250,
	}
	billRepo.Save(context.Background(), bill)

	err := billService.PayBill(context.Background(), bill.ID, 250)

	require.NoError(t, err, "failed to pay bill")
	bill, err = billRepo.FindByID(context.Background(), bill.ID)
	require.NoError(t, err, "failed to find bill")
	assert.Equal(t, domain.BillStatusPaid, bill.Status, "bill status not paid")
	assert.Equal(t, 250, bill.Paid, "bill already paid not correct")
}

func TestPayBillPartiallySuccess(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)
	bill := domain.Bill{
		ID:      id.New(),
		TableID: id.New(),
		Status:  domain.BillStatusPending,
		Items: []domain.MenuItem{
			{ID: id.New(), Name: "Spaghetti", Price: 100},
			{ID: id.New(), Name: "Pizza", Price: 150},
		},
		TotalAmount: 250,
	}
	billRepo.Save(context.Background(), bill)

	err := billService.PayBill(context.Background(), bill.ID, 100)

	require.NoError(t, err, "failed to pay bill")
	bill, err = billRepo.FindByID(context.Background(), bill.ID)
	require.NoError(t, err, "failed to find bill")
	assert.Equal(t, domain.BillPartiallyPaid, bill.Status, "bill status not partially paid")
	assert.Equal(t, 100, bill.Paid, "bill already paid not correct")
}

func TestPayBillAlreadyPaid(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)
	bill := domain.Bill{
		ID:          id.New(),
		TableID:     id.New(),
		Status:      domain.BillStatusPaid,
		TotalAmount: 250,
		Paid:        250,
	}
	billRepo.Save(context.Background(), bill)

	err := billService.PayBill(context.Background(), bill.ID, 100)

	require.Error(t, err, "paying already paid bill should fail")
}

func TestPayBillMoreThanTotalAmount(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)
	bill := domain.Bill{
		ID:          id.New(),
		TableID:     id.New(),
		Status:      domain.BillStatusPending,
		TotalAmount: 250,
	}
	billRepo.Save(context.Background(), bill)

	err := billService.PayBill(context.Background(), bill.ID, 300)

	require.Error(t, err, "paying more than total amount should fail")
}

func TestPayBillPartiallyMoreThanTotalAmount(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)
	bill := domain.Bill{
		ID:          id.New(),
		TableID:     id.New(),
		Status:      domain.BillStatusPending,
		TotalAmount: 200,
		Paid:        100,
	}
	billRepo.Save(context.Background(), bill)

	err := billService.PayBill(context.Background(), bill.ID, 101)

	require.Error(t, err, "paying more than total amount should fail")
}

func TestPayBillNotFound(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)

	err := billService.PayBill(context.Background(), id.New(), 100)

	require.Error(t, err, "paying not found bill should fail")
}
