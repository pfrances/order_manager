package usecases_test

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTakeOrderSuccess(t *testing.T) {
	orderRepo := memory.NewOrderRepository()
	kitchenRepo := memory.NewKitchenRepository()
	takeOrder := usecases.NewTakeOrder(orderRepo, kitchenRepo)
	tableID := id.NewID()
	menuItemIDs := []id.ID{id.NewID(), id.NewID()}

	orderID, err := takeOrder.Execute(tableID, menuItemIDs)

	require.NoError(t, err)
	order := orderRepo.GetOrder(orderID)
	assert.NotNil(t, order)
	assert.Equal(t, tableID, order.TableID)
	preparations := kitchenRepo.GetPreparationsByOrderID(order.ID)
	assert.Len(t, preparations, len(menuItemIDs))
	for _, preparation := range preparations {
		assert.Contains(t, menuItemIDs, preparation.MenuItemID)
		assert.Equal(t, order.ID, preparation.OrderID)
		assert.Equal(t, model.PreparationStatusPending, preparation.Status)
	}
}
