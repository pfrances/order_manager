package usecases_test

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTakeOrderSuccess(t *testing.T) {
	asserts := assert.New(t)
	orderRepo := memory.NewOrderRepository()
	kitchenRepo := memory.NewKitchenRepository()
	takeOrder := usecases.NewTakeOrder(orderRepo, kitchenRepo)

	tableID := id.NewID()
	menuItemIDs := []id.ID{id.NewID(), id.NewID()}

	orderID, err := takeOrder.Execute(tableID, menuItemIDs)
	asserts.NoError(err)
	asserts.NotEqual(id.NilID(), orderID)

	order := orderRepo.GetOrder(orderID)
	asserts.NotNil(order)
	asserts.Equal(tableID, order.TableID)

	preparations := kitchenRepo.GetPreparationsByOrderID(order.ID)
	asserts.Len(preparations, len(menuItemIDs))
	for _, preparation := range preparations {
		asserts.Equal(order.ID, preparation.OrderID)
		asserts.Contains(menuItemIDs, preparation.MenuItemID)
		asserts.Equal(model.PreparationStatusPending, preparation.Status)
	}
}
