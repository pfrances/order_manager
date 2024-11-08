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

func MustPresaveItemsFromTable(t *testing.T, db *sqlite.DB, table domain.Table) {
	t.Helper()

	items := make([]domain.MenuItem, 0)
	for _, order := range table.Orders {
		for _, preparation := range order.Preparations {
			items = append(items, preparation.MenuItem)
		}
	}

	itemRepo := sqlite.NewMenu(db)
	err := itemRepo.SaveItems(context.Background(), items)
	require.NoErrorf(t, err, "failed to save items: %v", err)
}

func GenerateDummyTable(status domain.TableStatus) domain.Table {
	menuItem1 := domain.MenuItem{ID: id.New(), Name: "item1", Price: 100}
	menuItem2 := domain.MenuItem{ID: id.New(), Name: "item2", Price: 200}
	menuItem3 := domain.MenuItem{ID: id.New(), Name: "item3", Price: 300}
	preparation1 := domain.Preparation{ID: id.New(), MenuItem: menuItem1, Status: domain.PreparationStatusPending}
	preparation2 := domain.Preparation{ID: id.New(), MenuItem: menuItem2, Status: domain.PreparationStatusReady}
	preparation3 := domain.Preparation{ID: id.New(), MenuItem: menuItem3, Status: domain.PreparationStatusServed}
	order1 := domain.Order{ID: id.New(), Status: domain.OrderStatusDone, Preparations: []domain.Preparation{preparation1, preparation2}}
	order2 := domain.Order{ID: id.New(), Status: domain.OrderStatusTaken, Preparations: []domain.Preparation{preparation3}}
	table := domain.Table{
		ID:     id.New(),
		Status: status,
		Orders: []domain.Order{order1, order2},
	}

	return table
}

func TestSaveAndRetrieveTableByID(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	table := GenerateDummyTable(domain.TableStatusOpened)
	tableRepo := sqlite.NewTable(db)
	MustPresaveItemsFromTable(t, db, table)

	err := tableRepo.Save(context.Background(), table)
	require.NoErrorf(t, err, "failed to save table: %v", err)

	gotTable, err := tableRepo.FindByID(context.Background(), table.ID)
	require.NoErrorf(t, err, "failed to retrieve table: %v", err)

	assert.Equal(t, table, gotTable)
}

func TestSaveAndRetrieveTableByStatus(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	openedStatus := GenerateDummyTable(domain.TableStatusOpened)
	closedStatus := GenerateDummyTable(domain.TableStatusClosed)
	tableRepo := sqlite.NewTable(db)
	MustPresaveItemsFromTable(t, db, openedStatus)
	MustPresaveItemsFromTable(t, db, closedStatus)

	err := tableRepo.Save(context.Background(), openedStatus)
	require.NoErrorf(t, err, "failed to save table: %v", err)

	err = tableRepo.Save(context.Background(), closedStatus)
	require.NoErrorf(t, err, "failed to save table: %v", err)

	gotOpenedTables, err := tableRepo.FindByStatus(context.Background(), domain.TableStatusOpened)
	require.NoErrorf(t, err, "failed to retrieve table: %v", err)

	assert.Len(t, gotOpenedTables, 1)
	assert.Equal(t, openedStatus, gotOpenedTables[0])

	gotClosedTables, err := tableRepo.FindByStatus(context.Background(), domain.TableStatusClosed)
	require.NoErrorf(t, err, "failed to retrieve table: %v", err)

	assert.Len(t, gotClosedTables, 1)
	assert.Equal(t, closedStatus, gotClosedTables[0])
}

func TestNotFoundTable(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	tableRepo := sqlite.NewTable(db)

	_, err := tableRepo.FindByID(context.Background(), id.New())
	assert.Equal(t, domain.ErrorCode(err), domain.ENOTFOUND)
}

func TestSaveTableWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	table := GenerateDummyTable(domain.TableStatusOpened)
	tableRepo := sqlite.NewTable(db)
	MustPresaveItemsFromTable(t, db, table)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := tableRepo.Save(ctx, table)
	assert.Error(t, err)

	_, err = tableRepo.FindByID(context.Background(), table.ID)
	assert.Error(t, err)

	assert.Equal(t, domain.ErrorCode(err), domain.ENOTFOUND)
}

func TestFindTableByStatusWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	table := GenerateDummyTable(domain.TableStatusOpened)
	tableRepo := sqlite.NewTable(db)
	MustPresaveItemsFromTable(t, db, table)

	err := tableRepo.Save(context.Background(), table)
	require.NoErrorf(t, err, "failed to save table: %v", err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = tableRepo.FindByStatus(ctx, domain.TableStatusOpened)
	assert.Error(t, err)
}
