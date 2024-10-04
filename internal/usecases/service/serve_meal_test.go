package service_test

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func orderWithPreparations(orderStatus model.OrderStatus, preparationsStatus []model.PreparationStatus) (model.Order, []model.Preparation) {
	menuItemIds := make([]id.ID, len(preparationsStatus))
	for i := range menuItemIds {
		menuItemIds[i] = id.NewID()
	}

	order := model.Order{
		ID:          id.NewID(),
		TableID:     id.NewID(),
		MenuItemIDs: menuItemIds,
		Status:      orderStatus,
	}

	preparations := make([]model.Preparation, len(order.MenuItemIDs))
	for i, status := range preparationsStatus {
		preparations[i] = model.Preparation{
			ID:         id.NewID(),
			OrderID:    order.ID,
			MenuItemID: order.MenuItemIDs[i],
			Status:     status,
		}
	}
	return order, preparations
}

func TestServeOrderSuccess(t *testing.T) {
	kitchenRepo := memory.NewKitchenRepository()
	orderRepo := memory.NewOrderRepository()
	usecase := service.NewServeMeal(orderRepo, kitchenRepo)

	tc := []struct {
		name                   string
		otherPreparationStatus model.PreparationStatus
		expectedOrderStatus    model.OrderStatus
	}{
		{
			name:                   "other preparation ready",
			otherPreparationStatus: model.PreparationStatusReady,
			expectedOrderStatus:    model.OrderStatusTaken,
		},
		{
			name:                   "other preparation served",
			otherPreparationStatus: model.PreparationStatusServed,
			expectedOrderStatus:    model.OrderStatusDone,
		},
		{
			name:                   "other preparation aborted",
			otherPreparationStatus: model.PreparationStatusAborted,
			expectedOrderStatus:    model.OrderStatusDone,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			order, preparations := orderWithPreparations(model.OrderStatusTaken, []model.PreparationStatus{model.PreparationStatusReady, tt.otherPreparationStatus})
			err := orderRepo.CreateOrder(order)
			require.NoError(t, err)
			err = kitchenRepo.CreatePreparations(preparations)
			require.NoError(t, err)
			preparationToServeID := preparations[0].ID

			err = usecase.Execute(preparationToServeID)

			require.NoError(t, err)
			updatedPreparation := kitchenRepo.GetPreparation(preparationToServeID)
			assert.Equal(t, model.PreparationStatusServed, updatedPreparation.Status)
			updatedOrder := orderRepo.GetOrder(order.ID)
			assert.Equal(t, tt.expectedOrderStatus, updatedOrder.Status)
		})
	}
}

func TestServeOrderFailWhenWrongStatus(t *testing.T) {
	kitchenRepo := memory.NewKitchenRepository()
	orderRepo := memory.NewOrderRepository()
	usecase := service.NewServeMeal(orderRepo, kitchenRepo)

	tc := []model.PreparationStatus{
		model.PreparationStatusPending,
		model.PreparationStatusInProgress,
		model.PreparationStatusServed,
		model.PreparationStatusAborted,
	}

	for _, tt := range tc {
		t.Run(string(tt), func(t *testing.T) {
			order, preparations := orderWithPreparations(model.OrderStatusTaken, []model.PreparationStatus{tt})
			err := orderRepo.CreateOrder(order)
			require.NoError(t, err)
			err = kitchenRepo.CreatePreparations(preparations)
			require.NoError(t, err)
			preparationToServeID := preparations[0].ID

			err = usecase.Execute(preparationToServeID)

			require.ErrorIs(t, err, model.ErrPreparationWrongStatus)
		})
	}
}

func TestServeOrderFailOrderNotFound(t *testing.T) {
	asserts := assert.New(t)
	kitchenRepo := memory.NewKitchenRepository()
	orderRepo := memory.NewOrderRepository()
	usecase := service.NewServeMeal(orderRepo, kitchenRepo)

	preparationToServeID := id.NewID()

	err := usecase.Execute(preparationToServeID)
	asserts.Error(model.ErrPreparationNotFound, err)
}
