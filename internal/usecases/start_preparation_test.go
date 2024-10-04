package usecases_test

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStartPreparationSuccess(t *testing.T) {
	repo := memory.NewKitchenRepository()
	startPreparation := usecases.NewStartPreparation(repo)
	preparation := model.Preparation{
		ID:      id.NewID(),
		OrderID: id.NewID(),
		Status:  model.PreparationStatusPending,
	}
	err := repo.CreatePreparation(preparation)
	require.NoError(t, err)

	err = startPreparation.Execute(preparation.ID)

	require.NoError(t, err)
	updatedPreparation := repo.GetPreparation(preparation.ID)
	require.Equal(t, model.PreparationStatusInProgress, updatedPreparation.Status)
}

func TestStartPreparationFailedWithWrongStatus(t *testing.T) {
	testCases := []model.PreparationStatus{
		model.PreparationStatusInProgress,
		model.PreparationStatusReady,
		model.PreparationStatusAborted,
	}
	repo := memory.NewKitchenRepository()
	startPreparation := usecases.NewStartPreparation(repo)

	for _, tc := range testCases {
		preparation := model.Preparation{
			ID:      id.NewID(),
			OrderID: id.NewID(),
			Status:  tc,
		}
		err := repo.CreatePreparation(preparation)
		require.NoError(t, err)

		err = startPreparation.Execute(preparation.ID)

		require.ErrorIs(t, err, model.ErrPreparationWrongStatus)
	}
}
