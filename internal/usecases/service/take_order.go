package service

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type TakeOrder struct {
	orderRepository   repositories.OrderRepository
	kitchenRepository repositories.KitchenRepository
}

func NewTakeOrder(orderRepository repositories.OrderRepository, kitchenRepository repositories.KitchenRepository) *TakeOrder {
	return &TakeOrder{orderRepository: orderRepository, kitchenRepository: kitchenRepository}
}

func (t *TakeOrder) Execute(tableId id.ID, menuItemIds []id.ID) (id.ID, error) {
	newOrder := model.Order{
		ID:          id.NewID(),
		TableID:     tableId,
		MenuItemIDs: menuItemIds,
		Status:      model.OrderStatusTaken,
	}

	err := t.orderRepository.CreateOrder(newOrder)
	if err != nil {
		return id.NilID(), fmt.Errorf("failed to create order %v: %w", newOrder, err)
	}

	preparations := make([]model.Preparation, 0, len(menuItemIds))
	for _, menuItemID := range menuItemIds {
		preparations = append(preparations, model.Preparation{
			ID:         id.NewID(),
			OrderID:    newOrder.ID,
			MenuItemID: menuItemID,
			Status:     model.PreparationStatusPending,
		})
	}

	err = t.kitchenRepository.CreatePreparations(preparations)
	if err != nil {
		return id.NilID(), fmt.Errorf("failed to create preparations %v: %w", preparations, err)
	}

	t.orderRepository.UpdateOrder(newOrder.ID, func(order *model.Order) error {
		order.PreparationIDs = make([]id.ID, 0, len(preparations))
		for _, preparation := range preparations {
			order.PreparationIDs = append(order.PreparationIDs, preparation.ID)
		}

		return nil
	})

	return newOrder.ID, nil
}
