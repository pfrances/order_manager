package memory

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type KitchenRepository struct {
	preparations map[id.ID]Preparation
}

func NewKitchenRepository() *KitchenRepository {
	return &KitchenRepository{
		preparations: make(map[id.ID]Preparation),
	}
}

func (k *KitchenRepository) CreatePreparation(preparation *model.Preparation) error {
	k.preparations[preparation.ID] = newPreparation(preparation)
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

func (k *KitchenRepository) UpdatePreparation(id id.ID, fn func(preparation *model.Preparation) error) error {
	if _, ok := k.preparations[id]; !ok {
		return fmt.Errorf("preparation not found")
	}

	preparation := k.preparations[id].toModel()
	if err := fn(&preparation); err != nil {
		return err
	}

	k.preparations[id] = newPreparation(&preparation)
	return nil
}

func (k *KitchenRepository) RemovePreparation(id id.ID) error {
	if _, ok := k.preparations[id]; !ok {
		return fmt.Errorf("preparation not found")
	}

	delete(k.preparations, id)
	return nil
}
