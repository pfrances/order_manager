package usecases

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type ServeMeal struct {
	orderRepository   repositories.OrderRepository
	kitchenRepository repositories.KitchenRepository
}

func NewServeMeal(orderRepository repositories.OrderRepository, kitchenRepository repositories.KitchenRepository) *ServeMeal {
	return &ServeMeal{orderRepository: orderRepository, kitchenRepository: kitchenRepository}
}

func (s *ServeMeal) Execute(preparationID id.ID) error {
	return s.kitchenRepository.UpdatePreparation(preparationID, func(preparation *model.Preparation) error {
		if preparation.Status != model.PreparationStatusReady {
			return model.ErrPreparationNotReady
		}

		preparation.Status = model.PreparationStatusServed

		preparations := s.kitchenRepository.GetPreparationsByOrderID(preparation.OrderID)
		allServed := true
		for _, preparation := range preparations {
			if preparation.ID == preparationID {
				continue
			}

			if preparation.Status != model.PreparationStatusServed && preparation.Status != model.PreparationStatusAborted {
				allServed = false
				break
			}
		}

		if allServed {
			return s.orderRepository.UpdateOrder(preparation.OrderID, func(order *model.Order) error {
				order.Status = model.OrderStatusDone
				return nil
			})
		}

		return nil
	})
}
