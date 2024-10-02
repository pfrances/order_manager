package usecases_test

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
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
	asserts := assert.New(t)
	kitchenRepo := memory.NewKitchenRepository()
	orderRepo := memory.NewOrderRepository()
	usecase := usecases.NewServeMeal(orderRepo, kitchenRepo)

	tc := []struct {
		name                   string
		otherPreparationStatus model.PreparationStatus
		expectedOrderStatus    model.OrderStatus
	}{
		{
			name:                   "other preparations ready",
			otherPreparationStatus: model.PreparationStatusReady,
			expectedOrderStatus:    model.OrderStatusTaken,
		},
		{
			name:                   "other preparations served",
			otherPreparationStatus: model.PreparationStatusServed,
			expectedOrderStatus:    model.OrderStatusDone,
		},
		{
			name:                   "other preparations aborted",
			otherPreparationStatus: model.PreparationStatusAborted,
			expectedOrderStatus:    model.OrderStatusDone,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			order, preparations := orderWithPreparations(model.OrderStatusTaken, []model.PreparationStatus{model.PreparationStatusReady, tt.otherPreparationStatus})
			err := orderRepo.CreateOrder(order)
			asserts.NoError(err)
			err = kitchenRepo.CreatePreparations(preparations)
			asserts.NoError(err)
			preparationToServeID := preparations[0].ID

			err = usecase.Execute(preparationToServeID)
			asserts.NoError(err)

			updatedPreparation := kitchenRepo.GetPreparation(preparationToServeID)
			asserts.Equal(model.PreparationStatusServed, updatedPreparation.Status)

			updatedOrder := orderRepo.GetOrder(order.ID)
			asserts.Equal(tt.expectedOrderStatus, updatedOrder.Status)
		})
	}
}

func TestServeOrderFailWhenWrongStatus(t *testing.T) {
	asserts := assert.New(t)
	kitchenRepo := memory.NewKitchenRepository()
	orderRepo := memory.NewOrderRepository()
	usecase := usecases.NewServeMeal(orderRepo, kitchenRepo)

	tc := []struct {
		name              string
		preparationStatus model.PreparationStatus
		expectedError     error
	}{
		{
			name:              "preparation pending",
			preparationStatus: model.PreparationStatusPending,
			expectedError:     model.ErrPreparationNotReady,
		},
		{
			name:              "preparation in progress",
			preparationStatus: model.PreparationStatusInProgress,
			expectedError:     model.ErrPreparationNotReady,
		},
		{
			name:              "preparation already served",
			preparationStatus: model.PreparationStatusServed,
			expectedError:     model.ErrPreparationNotReady,
		},
		{
			name:              "preparation aborted",
			preparationStatus: model.PreparationStatusAborted,
			expectedError:     model.ErrPreparationNotReady,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			order, preparations := orderWithPreparations(model.OrderStatusTaken, []model.PreparationStatus{tt.preparationStatus})
			err := orderRepo.CreateOrder(order)
			asserts.NoError(err)
			err = kitchenRepo.CreatePreparations(preparations)
			asserts.NoError(err)
			preparationToServeID := preparations[0].ID

			err = usecase.Execute(preparationToServeID)
			asserts.Error(tt.expectedError, err)
		})
	}
}

func TestServeOrderFailOrderNotFound(t *testing.T) {
	asserts := assert.New(t)
	kitchenRepo := memory.NewKitchenRepository()
	orderRepo := memory.NewOrderRepository()
	usecase := usecases.NewServeMeal(orderRepo, kitchenRepo)

	preparationToServeID := id.NewID()

	err := usecase.Execute(preparationToServeID)
	asserts.Error(model.ErrPreparationNotFound, err)
}
