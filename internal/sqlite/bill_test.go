package sqlite_test

import (
	"context"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"order_manager/internal/sqlite"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func MustPresaveTableFromBill(t *testing.T, db *sqlite.DB, bill domain.Bill) {
	t.Helper()

	items := bill.Items
	itemRepo := sqlite.NewMenu(db)

	err := itemRepo.SaveItems(context.Background(), items)
	require.NoErrorf(t, err, "failed to save items: %v", err)

	tableRepo := sqlite.NewTable(db)
	table := domain.Table{
		ID:     bill.TableID,
		Status: domain.TableStatusOpened,
		Orders: []domain.Order{
			{ID: id.New(), Status: domain.OrderStatusDone, Preparations: make([]domain.Preparation, 0, len(bill.Items))},
		},
	}

	for _, item := range bill.Items {
		table.Orders[0].Preparations = append(table.Orders[0].Preparations, domain.Preparation{
			ID:       id.New(),
			MenuItem: item,
			Status:   domain.PreparationStatusServed,
		})
	}

	err = tableRepo.Save(context.Background(), table)
	require.NoErrorf(t, err, "failed to save table: %v", err)
}

func GenerateDummyBill() domain.Bill {
	return domain.Bill{
		ID:          id.New(),
		TableID:     id.New(),
		TotalAmount: 300,
		Paid:        0,
		Status:      domain.BillStatusPending,
		Items: []domain.MenuItem{
			{
				ID:    id.New(),
				Name:  "item1",
				Price: 100,
			},
			{
				ID:    id.New(),
				Name:  "item2",
				Price: 200,
			},
		},
	}
}

func TestSaveAndRetriveBillByID(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	bill := GenerateDummyBill()
	billRepo := sqlite.NewBill(db)
	MustPresaveTableFromBill(t, db, bill)

	err := billRepo.Save(context.Background(), bill)
	require.NoErrorf(t, err, "failed to save bill: %v", err)

	gotBill, err := billRepo.FindByID(context.Background(), bill.ID)
	require.NoErrorf(t, err, "failed to retrieve bill: %v", err)

	assert.Equal(t, bill, gotBill)
}

func TestSaveAndRetrieveBillByTableID(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	bill := GenerateDummyBill()
	billRepo := sqlite.NewBill(db)
	MustPresaveTableFromBill(t, db, bill)

	err := billRepo.Save(context.Background(), bill)
	require.NoErrorf(t, err, "failed to save bill: %v", err)

	gotBills, err := billRepo.FindByTableID(context.Background(), bill.TableID)
	require.NoErrorf(t, err, "failed to retrieve bill: %v", err)

	assert.Equal(t, []domain.Bill{bill}, gotBills)
}

func TestNotFoundBillByID(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	billRepo := sqlite.NewBill(db)

	_, err := billRepo.FindByID(context.Background(), id.New())
	assert.Equal(t, domain.ErrorCode(err), domain.ENOTFOUND)
}

func TestSaveBillWithInvalidTableID(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	bill := GenerateDummyBill()
	bill.TableID = id.New()
	billRepo := sqlite.NewBill(db)

	err := billRepo.Save(context.Background(), bill)
	assert.Error(t, err)
}

func TestSaveBillWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	bill := GenerateDummyBill()
	billRepo := sqlite.NewBill(db)
	MustPresaveTableFromBill(t, db, bill)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := billRepo.Save(ctx, bill)
	assert.Error(t, err)
}

func TestFindBillWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	billRepo := sqlite.NewBill(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := billRepo.FindByID(ctx, id.New())
	assert.Error(t, err)
}

func TestFindBillsByTableIDWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	billRepo := sqlite.NewBill(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := billRepo.FindByTableID(ctx, id.New())
	assert.Error(t, err)
}
