package usecases

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type StartPreparation struct {
	KitchenRepository repositories.KitchenRepository
}

func NewStartPreparation(kitchenRepository repositories.KitchenRepository) *StartPreparation {
	return &StartPreparation{KitchenRepository: kitchenRepository}
}

func (s *StartPreparation) Execute(orderId id.ID) error {
	return s.KitchenRepository.UpdatePreparation(orderId, func(preparation *model.Preparation) error {
		if preparation.Status != model.PreparationStatusPending {
			return fmt.Errorf("%w: start preparation failed. preparation with ID %s has status %s, but only %s is allowed",
				model.ErrPreparationWrongStatus, orderId, preparation.Status, model.PreparationStatusPending)
		}

		preparation.Status = model.PreparationStatusInProgress
		return nil
	})
}
