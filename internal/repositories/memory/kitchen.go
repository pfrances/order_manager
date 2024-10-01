package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type KitchenRepository struct {
	preparations map[id.ID]Preparation
}

func NewKitchenRepository() *KitchenRepository {
	repo := &KitchenRepository{
		preparations: make(map[id.ID]Preparation),
	}

	return repo
}

func (k *KitchenRepository) CreatePreparation(preparation model.Preparation) error {
	if _, ok := k.preparations[preparation.ID]; ok {
		return repositories.ErrPreparationAlreadyExists
	}

	k.preparations[preparation.ID] = fromModel(preparation)
	return nil
}

func (k *KitchenRepository) CreatePreparations(preparations []model.Preparation) error {
	for _, p := range preparations {
		if err := k.CreatePreparation(p); err != nil {
			return err
		}
	}

	return nil
}

func (k *KitchenRepository) GetPreparation(id id.ID) *model.Preparation {
	preparation, ok := k.preparations[id]
	if !ok {
		return nil
	}

	modelPreparation := preparation.toModel()
	return &modelPreparation
}

func (k *KitchenRepository) GetPreparationsByOrderID(orderID id.ID) []model.Preparation {
	var preparations []model.Preparation
	for _, preparation := range k.preparations {
		if preparation.orderID == orderID {
			preparations = append(preparations, preparation.toModel())
		}
	}

	return preparations
}

func (k *KitchenRepository) UpdatePreparation(id id.ID, fn func(preparation *model.Preparation) error) error {
	preparation := k.GetPreparation(id)
	if preparation == nil {
		return repositories.ErrPreparationNotFound
	}

	if err := fn(preparation); err != nil {
		return err
	}

	k.preparations[id] = fromModel(*preparation)
	return nil
}

func (k *KitchenRepository) RemovePreparation(id id.ID) error {
	if _, ok := k.preparations[id]; !ok {
		return repositories.ErrPreparationNotFound
	}

	delete(k.preparations, id)
	return nil
}
