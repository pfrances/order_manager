package usecases_test

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			expectedError: model.ErrPreparationNotPending,
		},
		{
			name:          "ready",
			status:        model.PreparationStatusReady,
			expectedError: model.ErrPreparationNotPending,
		},
		{
			name:          "aborted",
			status:        model.PreparationStatusAborted,
			expectedError: model.ErrPreparationNotPending,
		},
	}

	asserts := assert.New(t)
	repo := memory.NewKitchenRepository()
	startPreparation := usecases.NewStartPreparation(repo)

	for _, tc := range testCases {
		preparation := model.Preparation{
			ID:      id.NewID(),
			OrderID: id.NewID(),
			Status:  tc.status,
		}
		repo.CreatePreparation(preparation)

		err := startPreparation.Execute(preparation.ID)

		asserts.Equal(tc.expectedError, err, tc.name)
	}
}
