package kitchen

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type FinishPreparation struct {
	kitchenRepository repositories.KitchenRepository
}

func NewFinishPreparation(kitchenRepo repositories.KitchenRepository) *FinishPreparation {
	return &FinishPreparation{kitchenRepository: kitchenRepo}
}

func (f *FinishPreparation) Execute(PreparationID id.ID) error {
	return f.kitchenRepository.UpdatePreparation(PreparationID, func(preparation *model.Preparation) error {
		if preparation.Status != model.PreparationStatusInProgress {
			return fmt.Errorf("finish preparation failed. preparation with ID %s has status %s, but only %s is allowed: %w",
				PreparationID, preparation.Status, model.PreparationStatusInProgress, model.ErrPreparationWrongStatus)
		}

		preparation.Status = model.PreparationStatusReady
		return nil
	})
}
