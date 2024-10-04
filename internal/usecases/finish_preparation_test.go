package usecases_test

import (
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinishPreparationSuccess(t *testing.T) {
	kitchenRepo := memory.NewKitchenRepository()
	usecase := usecases.NewFinishPreparation(kitchenRepo)

	_, preparations := orderWithPreparations(model.OrderStatusTaken, []model.PreparationStatus{model.PreparationStatusInProgress})
	kitchenRepo.CreatePreparation(preparations[0])

	err := usecase.Execute(preparations[0].ID)
	require.NoError(t, err)

	preparation := kitchenRepo.GetPreparation(preparations[0].ID)
	assert.Equal(t, model.PreparationStatusReady, preparation.Status)
}

func TestFinishPreparationWhenWrongStatus(t *testing.T) {
	kitchenRepo := memory.NewKitchenRepository()
	usecase := usecases.NewFinishPreparation(kitchenRepo)

	tt := []model.PreparationStatus{
		model.PreparationStatusReady,
		model.PreparationStatusServed,
		model.PreparationStatusAborted,
	}
	for _, status := range tt {
		t.Run(string(status), func(t *testing.T) {
			_, preparations := orderWithPreparations(model.OrderStatusTaken, []model.PreparationStatus{status})
			kitchenRepo.CreatePreparation(preparations[0])

			err := usecase.Execute(preparations[0].ID)

			require.ErrorIs(t, err, model.ErrPreparationWrongStatus)
		})

	}

}
