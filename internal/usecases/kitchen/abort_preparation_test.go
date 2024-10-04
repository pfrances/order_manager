package kitchen_test

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases/kitchen"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAbordPreparationSuccess(t *testing.T) {
	repo := memory.NewKitchenRepository()
	abortPreparation := kitchen.NewAbortPreparation(repo)

	tt := []model.PreparationStatus{
		model.PreparationStatusPending,
		model.PreparationStatusInProgress,
	}
	for _, status := range tt {
		preparation := model.Preparation{
			ID:      id.NewID(),
			OrderID: id.NewID(),
			Status:  status,
		}
		err := repo.CreatePreparation(preparation)
		require.NoError(t, err)

		err = abortPreparation.Execute(preparation.ID)

		require.NoError(t, err)
		updatedPreparation := repo.GetPreparation(preparation.ID)
		assert.Equal(t, model.PreparationStatusAborted, updatedPreparation.Status)
	}
}

func TestAbordPreparationFailWenWrongStatus(t *testing.T) {
	repo := memory.NewKitchenRepository()
	abortPreparation := kitchen.NewAbortPreparation(repo)

	tt := []model.PreparationStatus{
		model.PreparationStatusReady,
		model.PreparationStatusServed,
		model.PreparationStatusAborted,
	}
	for _, status := range tt {
		preparation := model.Preparation{
			ID:      id.NewID(),
			OrderID: id.NewID(),
			Status:  status,
		}
		err := repo.CreatePreparation(preparation)
		require.NoError(t, err)

		err = abortPreparation.Execute(preparation.ID)

		require.ErrorIs(t, err, model.ErrPreparationWrongStatus)
		notUpdatedPreparation := repo.GetPreparation(preparation.ID)
		assert.Equal(t, status, notUpdatedPreparation.Status)
	}
}
