package service_test

import (
	"context"
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTakeOrderSuccess(t *testing.T) {
	txManager := memory.NewTransactionManager()
	tableRepo := memory.NewTableRepository()
	orderRepo := memory.NewOrderRepository()
	kitchenRepo := memory.NewKitchenRepository()
	takeOrder := service.NewTakeOrder(txManager, tableRepo, orderRepo, kitchenRepo)

	table := model.Table{ID: id.NewID()}
	tableRepo.CreateTable(table)
	menuItemIDs := []id.ID{id.NewID(), id.NewID()}

	orderID, err := takeOrder.Execute(context.Background(), table.ID, menuItemIDs)
	require.NoError(t, err)

	order := orderRepo.GetOrder(orderID)
	assert.Equal(t, table.ID, order.TableID)

	updatedTable := tableRepo.GetTable(table.ID)
	assert.Contains(t, updatedTable.OrderIDs, orderID)

	preparations := kitchenRepo.GetPreparationsByOrderID(order.ID)
	assert.Len(t, preparations, len(menuItemIDs))
	for _, preparation := range preparations {
		assert.Contains(t, menuItemIDs, preparation.MenuItemID)
		assert.Equal(t, order.ID, preparation.OrderID)
		assert.Equal(t, model.PreparationStatusPending, preparation.Status)
	}
}
