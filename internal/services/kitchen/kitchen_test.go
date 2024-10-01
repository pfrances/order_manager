package kitchen_test

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/services/kitchen"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePreparation(t *testing.T) {
	kitchenRepo := memory.NewKitchenRepository()
	kitchen := kitchen.NewService(kitchenRepo)

	order := model.Order{
		ID: id.NewID(),
	}

	_, err := kitchen.CreatePreparation(order.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetPreparation(t *testing.T) {
	asserts := assert.New(t)
	repo := memory.NewKitchenRepository()
	service := kitchen.NewService(repo)

	order := model.Order{
		ID: id.NewID(),
	}

	preparationID, err := service.CreatePreparation(order.ID)
	asserts.Nil(err)

	preparation := service.GetPreparation(preparationID)
	asserts.NotNil(preparation)
}

func TestStartPreparation(t *testing.T) {
	testCases := []struct {
		name          string
		status        model.PreparationStatus
		expectedError error
	}{
		{
			name:          "pending",
			status:        model.PreparationStatusPending,
			expectedError: nil,
		},
		{
			name:          "in progress",
			status:        model.PreparationStatusInProgress,
			expectedError: kitchen.ErrPreparationNotPending,
		},
		{
			name:          "ready",
			status:        model.PreparationStatusReady,
			expectedError: kitchen.ErrPreparationNotPending,
		},
		{
			name:          "aborted",
			status:        model.PreparationStatusAborted,
			expectedError: kitchen.ErrPreparationNotPending,
		},
	}

	asserts := assert.New(t)
	repo := memory.NewKitchenRepository()
	service := kitchen.NewService(repo)

	for _, tc := range testCases {
		preparation := model.Preparation{
			ID:      id.NewID(),
			OrderID: id.NewID(),
			Status:  tc.status,
		}
		repo.CreatePreparation(preparation)

		err := service.StartPreparation(preparation.ID)

		asserts.Equal(tc.expectedError, err, tc.name)
	}
}

func TestFinishPreparation(t *testing.T) {
	testCases := []struct {
		name          string
		status        model.PreparationStatus
		expectedError error
	}{
		{
			name:          "pending",
			status:        model.PreparationStatusPending,
			expectedError: kitchen.ErrPreparationNotInProgress,
		},
		{
			name:          "in progress",
			status:        model.PreparationStatusInProgress,
			expectedError: nil,
		},
		{
			name:          "ready",
			status:        model.PreparationStatusReady,
			expectedError: kitchen.ErrPreparationNotInProgress,
		},
		{
			name:          "aborted",
			status:        model.PreparationStatusAborted,
			expectedError: kitchen.ErrPreparationNotInProgress,
		},
	}

	asserts := assert.New(t)
	repo := memory.NewKitchenRepository()
	service := kitchen.NewService(repo)

	for _, tc := range testCases {
		preparation := model.Preparation{
			ID:      id.NewID(),
			OrderID: id.NewID(),
			Status:  tc.status,
		}
		repo.CreatePreparation(preparation)

		err := service.FinishPreparation(preparation.ID)

		asserts.Equal(tc.expectedError, err, tc.name)
	}
}

func TestAbortPreparation(t *testing.T) {
	testCases := []struct {
		name          string
		status        model.PreparationStatus
		expectedError error
	}{
		{
			name:          "pending",
			status:        model.PreparationStatusPending,
			expectedError: nil,
		},
		{
			name:          "in progress",
			status:        model.PreparationStatusInProgress,
			expectedError: nil,
		},
		{
			name:          "ready",
			status:        model.PreparationStatusReady,
			expectedError: nil,
		},
		{
			name:          "aborted",
			status:        model.PreparationStatusAborted,
			expectedError: kitchen.ErrPreparationAlreadyAborted,
		},
	}

	asserts := assert.New(t)
	repo := memory.NewKitchenRepository()
	service := kitchen.NewService(repo)

	for _, tc := range testCases {
		preparation := model.Preparation{
			ID:      id.NewID(),
			OrderID: id.NewID(),
			Status:  tc.status,
		}
		repo.CreatePreparation(preparation)

		err := service.AbortPreparation(preparation.ID)

		asserts.Equal(tc.expectedError, err, tc.name)
	}
}
