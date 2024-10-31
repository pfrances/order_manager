package bill_test

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/money"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases/bill"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBillGenerationSuccess(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	tableRepo := memory.NewTableRepository()
	tableID := id.NewID()
	item1 := model.MenuItem{ID: id.NewID(), Name: "Item 1", Price: 10}
	item2 := model.MenuItem{ID: id.NewID(), Name: "Item 2", Price: 20}
	menuRepo := memory.NewMenuRepository()
	menuRepo.CreateMenuItem(item1)
	menuRepo.CreateMenuItem(item2)
	order := model.Order{ID: id.NewID(), TableID: tableID, MenuItemIDs: []id.ID{item1.ID, item2.ID}}
	orderRepo.CreateOrder(order)
	table := model.Table{ID: id.NewID()}
	tableRepo.CreateTable(table)
	tableRepo.UpdateTable(table.ID, func(table *model.Table) error {
		table.OrderIDs = append(table.OrderIDs, order.ID)
		return nil
	})
	usecase := bill.NewGenerateBill(menuRepo, orderRepo, tableRepo)

	bill, err := usecase.Execute(table.ID)

	require.NoError(t, err)
	require.Equal(t, money.Money(30), bill.TotalAmount)
	require.Contains(t, bill.Items, item1)
	require.Contains(t, bill.Items, item2)
	require.Equal(t, money.Money(0), bill.AlreadyPaid)
}
