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

func TestBillIsValid(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			bill     domain.Bill
		}{
			{
				testName: "valid pending bill",
				bill: domain.Bill{
					ID:          id.New(),
					TableID:     id.New(),
					Items:       []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
					Status:      domain.BillStatusPending,
					TotalAmount: 100,
					Paid:        0,
				},
			},
			{
				testName: "valid partially paid bill",
				bill: domain.Bill{
					ID:          id.New(),
					TableID:     id.New(),
					Items:       []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
					Status:      domain.BillPartiallyPaid,
					TotalAmount: 100,
					Paid:        50,
				},
			},
			{
				testName: "valid paid bill",
				bill: domain.Bill{
					ID:          id.New(),
					TableID:     id.New(),
					Items:       []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
					Status:      domain.BillStatusPaid,
					TotalAmount: 100,
					Paid:        100,
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				assert.Truef(t, tc.bill.IsValid(), "should be valid bill: %v", tc.bill)
			})
		}
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			bill     domain.Bill
		}{
			{
				testName: "invalid bill ID",
				bill: domain.Bill{
					ID:          id.NilID(),
					TableID:     id.New(),
					Items:       []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
					Status:      domain.BillStatusPending,
					TotalAmount: 100,
					Paid:        0,
				},
			},
			{
				testName: "invalid table ID",
				bill: domain.Bill{
					ID:          id.New(),
					TableID:     id.NilID(),
					Items:       []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
					Status:      domain.BillStatusPending,
					TotalAmount: 100,
					Paid:        0,
				},
			},
			{
				testName: "invalid item",
				bill: domain.Bill{
					ID:          id.New(),
					TableID:     id.New(),
					Items:       []domain.MenuItem{{ID: id.NilID(), Name: "test", Price: 100}},
					Status:      domain.BillStatusPending,
					TotalAmount: 100,
					Paid:        0,
				},
			},
			{
				testName: "invalid status",
				bill: domain.Bill{
					ID:          id.New(),
					TableID:     id.New(),
					Items:       []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
					Status:      "invalid",
					TotalAmount: 100,
					Paid:        0,
				},
			},
			{
				testName: "invalid total amount",
				bill: domain.Bill{
					ID:          id.New(),
					TableID:     id.New(),
					Items:       []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
					Status:      domain.BillStatusPending,
					TotalAmount: -100,
					Paid:        0,
				},
			},
			{
				testName: "invalid paid amount",
				bill: domain.Bill{
					ID:          id.New(),
					TableID:     id.New(),
					Items:       []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
					Status:      domain.BillStatusPending,
					TotalAmount: 100,
					Paid:        -100,
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				assert.Falsef(t, tc.bill.IsValid(), "should be invalid bill: %v", tc.bill)
			})
		}
	})
}

func TestCreateBill(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName       string
			table          domain.Table
			expectedAmount int
		}{
			{
				testName: "valid table with orders",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusClosed,
					Orders: []domain.Order{
						{ID: id.New(), Preparations: []domain.Preparation{
							{MenuItem: domain.MenuItem{ID: id.New(), Name: "Spaghetti", Price: 100}},
							{MenuItem: domain.MenuItem{ID: id.New(), Name: "Pizza", Price: 150}},
						}},
					},
				},
				expectedAmount: 250,
			},
			{
				testName: "valid table with no orders",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusClosed,
					Orders: make([]domain.Order, 0),
				},
				expectedAmount: 0,
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				bill, err := billService.GenerateBill(context.Background(), tc.table)

				require.NoError(t, err, "failed to generate bill")
				assert.Equal(t, domain.BillStatusPending, bill.Status, "bill status not pending")
				assert.Equal(t, tc.expectedAmount, bill.TotalAmount, "bill total amount not correct")
			})
		}
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			table    domain.Table
		}{
			{
				testName: "opened table",
				table: domain.Table{
					ID:     id.New(),
					Status: domain.TableStatusOpened,
					Orders: []domain.Order{
						{ID: id.New(), Preparations: []domain.Preparation{
							{MenuItem: domain.MenuItem{ID: id.New(), Name: "Spaghetti", Price: 100}},
							{MenuItem: domain.MenuItem{ID: id.New(), Name: "Pizza", Price: 150}},
						}},
					},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				_, err := billService.GenerateBill(context.Background(), tc.table)

				require.Error(t, err, "generating bill should fail")
			})
		}

		t.Run("table not found", func(t *testing.T) {
			t.Parallel()

			_, err := billService.GenerateBill(context.Background(), domain.Table{ID: id.New()})

			require.Error(t, err, "generating bill should fail")
		})

		t.Run("Context error", func(t *testing.T) {
			t.Parallel()

			_, err := billService.GenerateBill(context.Background(), domain.Table{ID: id.New()})

			require.Error(t, err, "generating bill should fail")
		})
	})
}

func TestPayBill(t *testing.T) {
	billRepo := inmem.NewBill()
	billService := domain.NewBillService(billRepo)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			bill     domain.Bill
			amount   int
		}{
			{
				testName: "pay full amount",
				bill: domain.Bill{
					ID:      id.New(),
					TableID: id.New(),
					Status:  domain.BillStatusPending,
					Items: []domain.MenuItem{
						{ID: id.New(), Name: "Spaghetti", Price: 100},
						{ID: id.New(), Name: "Pizza", Price: 150},
					},
					TotalAmount: 250,
				},
				amount: 250,
			},
			{
				testName: "pay partially",
				bill: domain.Bill{
					ID:      id.New(),
					TableID: id.New(),
					Status:  domain.BillStatusPending,
					Items: []domain.MenuItem{
						{ID: id.New(), Name: "Spaghetti", Price: 100},
						{ID: id.New(), Name: "Pizza", Price: 150},
					},
					TotalAmount: 250,
				},
				amount: 100,
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				err := billRepo.Save(context.Background(), tc.bill)
				require.NoError(t, err, "failed to save bill")

				err = billService.PayBill(context.Background(), tc.bill.ID, tc.amount)

				require.NoError(t, err, "failed to pay bill")
				bill, err := billRepo.FindByID(context.Background(), tc.bill.ID)
				require.NoError(t, err, "failed to find bill")
				assert.Equal(t, tc.amount, bill.Paid, "bill paid amount not correct")
			})
		}
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			bill     domain.Bill
			amount   int
		}{
			{
				testName: "pay already paid bill",
				bill: domain.Bill{
					ID:      id.New(),
					TableID: id.New(),
					Status:  domain.BillStatusPaid,
					Items: []domain.MenuItem{
						{ID: id.New(), Name: "Spaghetti", Price: 100},
						{ID: id.New(), Name: "Pizza", Price: 150},
					},
					TotalAmount: 250,
					Paid:        250,
				},
				amount: 100,
			},
			{
				testName: "pay more than total amount",
				bill: domain.Bill{
					ID:      id.New(),
					TableID: id.New(),
					Status:  domain.BillStatusPending,
					Items: []domain.MenuItem{
						{ID: id.New(), Name: "Spaghetti", Price: 100},
						{ID: id.New(), Name: "Pizza", Price: 150},
					},
					TotalAmount: 250,
				},
				amount: 300,
			},
			{
				testName: "pay partially more than total amount",
				bill: domain.Bill{
					ID:          id.New(),
					TableID:     id.New(),
					Status:      domain.BillStatusPending,
					Items:       []domain.MenuItem{{ID: id.New(), Name: "test", Price: 200}},
					TotalAmount: 200,
					Paid:        100,
				},
				amount: 101,
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()

				err := billRepo.Save(context.Background(), tc.bill)
				require.NoError(t, err, "failed to save bill")

				err = billService.PayBill(context.Background(), tc.bill.ID, tc.amount)

				require.Error(t, err, "paying bill should fail")
			})
		}

		t.Run("bill not found", func(t *testing.T) {
			t.Parallel()

			err := billService.PayBill(context.Background(), id.New(), 100)

			require.Error(t, err, "paying not found bill should fail")
		})

		t.Run("context error", func(t *testing.T) {
			t.Parallel()

			err := billService.PayBill(context.Background(), id.New(), 100)

			require.Error(t, err, "paying bill should fail")
		})
	})

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
