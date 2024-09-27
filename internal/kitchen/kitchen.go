package kitchen

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type Kitchen struct {
	repo Repository
}

type Repository interface {
	CreatePreparation(preparation *model.Preparation) error
	GetPreparation(id id.ID) *model.Preparation
	UpdatePreparation(id id.ID, fn func(preparation *model.Preparation) error) error
	RemovePreparation(id id.ID) error
}

func NewKitchen(repo Repository) *Kitchen {
	return &Kitchen{
		repo: repo,
	}
}

func (k *Kitchen) CreatePreparation(orderID id.ID) (id.ID, error) {
	id := id.NewID()

	preparation := &model.Preparation{
		ID:      id,
		OrderID: orderID,
		Status:  model.PreparationStatusPending,
	}
	return id, k.repo.CreatePreparation(preparation)
}

func (k *Kitchen) GetPreparation(id id.ID) *model.Preparation {
	return k.repo.GetPreparation(id)
}

func (k *Kitchen) UpdatePreparationStatus(id id.ID, status model.PreparationStatus) error {
	return k.repo.UpdatePreparation(id, func(preparation *model.Preparation) error {
		preparation.Status = status
		return nil
	})
}
