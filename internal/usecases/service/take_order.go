package service

import (
	"context"
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type TakeOrder struct {
	txManager         repositories.TransactionManager
	tableRepository   repositories.TableRepository
	orderRepository   repositories.OrderRepository
	kitchenRepository repositories.KitchenRepository
}

func NewTakeOrder(
	txManager repositories.TransactionManager,
	tableRepository repositories.TableRepository,
	orderRepository repositories.OrderRepository,
	kitchenRepository repositories.KitchenRepository,
) *TakeOrder {

	return &TakeOrder{
		txManager:         txManager,
		tableRepository:   tableRepository,
		orderRepository:   orderRepository,
		kitchenRepository: kitchenRepository,
	}
}

func (t *TakeOrder) Execute(ctx context.Context, tableID id.ID, menuItemIDs []id.ID) (id.ID, error) {
	txCtx, err := t.txManager.Begin(ctx)
	if err != nil {
		return id.NilID(), fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			t.txManager.Rollback(txCtx)
		} else {
			t.txManager.Commit(txCtx)
		}
	}()

	newOrder := model.Order{
		ID:          id.NewID(),
		TableID:     tableID,
		MenuItemIDs: menuItemIDs,
		Status:      model.OrderStatusTaken,
	}
	err = t.orderRepository.CreateOrder(newOrder)
	if err != nil {
		return id.NilID(), fmt.Errorf("failed to create order %v: %w", newOrder, err)
	}

	err = t.tableRepository.UpdateTable(tableID, func(table *model.Table) error {
		table.OrderIDs = append(table.OrderIDs, newOrder.ID)
		return nil
	})
	if err != nil {
		return id.NilID(), fmt.Errorf("failed to update table %s: %w", tableID, err)
	}

	preparations := make([]model.Preparation, 0, len(menuItemIDs))
	for _, menuItemID := range menuItemIDs {
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

	err = t.orderRepository.UpdateOrder(newOrder.ID, func(order *model.Order) error {
		order.PreparationIDs = make([]id.ID, 0, len(preparations))
		for _, preparation := range preparations {
			order.PreparationIDs = append(order.PreparationIDs, preparation.ID)
		}

		return nil
	})
	if err != nil {
		return id.NilID(), fmt.Errorf("failed to update order %s: %w", newOrder.ID, err)
	}

	return newOrder.ID, nil
}
