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

func TestFinishPreparationSuccess(t *testing.T) {
	kitchenRepo := memory.NewKitchenRepository()
	usecase := kitchen.NewFinishPreparation(kitchenRepo)
	preparation := model.Preparation{
		ID:      id.NewID(),
		OrderID: id.NewID(),
		Status:  model.PreparationStatusInProgress,
	}
	kitchenRepo.CreatePreparation(preparation)

	err := usecase.Execute(preparation.ID)
	require.NoError(t, err)

	UpdatePreparation := kitchenRepo.GetPreparation(preparation.ID)
	assert.Equal(t, model.PreparationStatusReady, UpdatePreparation.Status)
}

func TestFinishPreparationWhenWrongStatus(t *testing.T) {
	kitchenRepo := memory.NewKitchenRepository()
	usecase := kitchen.NewFinishPreparation(kitchenRepo)

	tt := []model.PreparationStatus{
		model.PreparationStatusReady,
		model.PreparationStatusServed,
		model.PreparationStatusAborted,
	}
	for _, status := range tt {
		t.Run(string(status), func(t *testing.T) {
			preparation := model.Preparation{
				ID:      id.NewID(),
				OrderID: id.NewID(),
				Status:  status,
			}
			kitchenRepo.CreatePreparation(preparation)

			err := usecase.Execute(preparation.ID)

			require.ErrorIs(t, err, model.ErrPreparationWrongStatus)
		})

	}

}
