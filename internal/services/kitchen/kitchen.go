package kitchen

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

var (
	ErrPreparationNotPending     = fmt.Errorf("preparation is not pending")
	ErrPreparationNotInProgress  = fmt.Errorf("preparation is not in progress")
	ErrPreparationAlreadyAborted = fmt.Errorf("preparation is already aborted")
)

type Service struct {
	repo repositories.KitchenRepository
}

func NewService(repo repositories.KitchenRepository) *Service {
	return &Service{repo: repo}
}

func (k *Service) CreatePreparation(orderID id.ID) (id.ID, error) {
	id := id.NewID()

	preparation := model.Preparation{
		ID:      id,
		OrderID: orderID,
		Status:  model.PreparationStatusPending,
	}
	return id, k.repo.CreatePreparation(preparation)
}

func (k *Service) GetPreparation(id id.ID) *model.Preparation {
	return k.repo.GetPreparation(id)
}

func (k *Service) StartPreparation(id id.ID) error {
	return k.repo.UpdatePreparation(id, func(preparation *model.Preparation) error {
		if preparation.Status != model.PreparationStatusPending {
			return ErrPreparationNotPending
		}

		preparation.Status = model.PreparationStatusInProgress
		return nil
	})
}

func (k *Service) FinishPreparation(id id.ID) error {
	return k.repo.UpdatePreparation(id, func(preparation *model.Preparation) error {
		if preparation.Status != model.PreparationStatusInProgress {
			return ErrPreparationNotInProgress
		}

		preparation.Status = model.PreparationStatusReady
		return nil
	})
}

func (k *Service) AbortPreparation(id id.ID) error {
	return k.repo.UpdatePreparation(id, func(preparation *model.Preparation) error {
		if preparation.Status == model.PreparationStatusAborted {
			return ErrPreparationAlreadyAborted
		}

		preparation.Status = model.PreparationStatusAborted
		return nil
	})
}
