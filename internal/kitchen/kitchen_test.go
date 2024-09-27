package kitchen_test

import (
	"order_manager/internal/id"
	"order_manager/internal/kitchen"
	"order_manager/internal/memory"
	"order_manager/internal/model"
	"testing"
)

func TestCreatePreparation(t *testing.T) {
	kitchenRepo := memory.NewKitchenRepository()
	kitchen := kitchen.NewKitchen(kitchenRepo)

	order := model.Order{
		ID: id.NewID(),
	}

	_, err := kitchen.CreatePreparation(order.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetPreparation(t *testing.T) {
	kitchenRepo := memory.NewKitchenRepository()
	kitchen := kitchen.NewKitchen(kitchenRepo)

	order := model.Order{
		ID: id.NewID(),
	}

	preparationID, err := kitchen.CreatePreparation(order.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	preparation := kitchen.GetPreparation(preparationID)
	if preparation == nil {
		t.Fatalf("expected preparation to be found")
	}
}
