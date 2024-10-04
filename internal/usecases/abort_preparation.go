package usecases

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type AbortPreparation struct {
	KitchenRepository repositories.KitchenRepository
}

func NewAbortPreparation(kitchenRepository repositories.KitchenRepository) *AbortPreparation {
	return &AbortPreparation{KitchenRepository: kitchenRepository}
}

func (a *AbortPreparation) Execute(preparationID id.ID) error {
	return a.KitchenRepository.UpdatePreparation(preparationID, func(preparation *model.Preparation) error {
		if preparation.Status != model.PreparationStatusPending && preparation.Status != model.PreparationStatusInProgress {
			return fmt.Errorf(
				"abort preparation failed. preparation with ID %s has status %s, but only %s and %s are allowed: %w",
				preparationID, preparation.Status, model.PreparationStatusPending, model.PreparationStatusInProgress, model.ErrPreparationWrongStatus)
		}

		preparation.Status = model.PreparationStatusAborted
		return nil
	})
}
