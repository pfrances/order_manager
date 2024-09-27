package menu_test

import (
	"order_manager/internal/memory"
	"order_manager/internal/menu"
	"testing"
)

func TestAddItem(t *testing.T) {
	menuRepo := memory.NewMenuRepository()
	menu := menu.NewMenuService(menuRepo)

	id, err := menu.CreateMenu("Cheeseburger", 10)
	if err != nil {
		t.Errorf("Error creating menu item: %v", err)
	}

	if menuRepo.GetMenuItem(id) == nil {
		t.Errorf("Menu item not found")
	}
}

func TestRemoveItem(t *testing.T) {
	menuRepo := memory.NewMenuRepository()
	menu := menu.NewMenuService(menuRepo)

	id, err := menu.CreateMenu("Cheeseburger", 10)
	if err != nil {
		t.Errorf("Error creating menu item: %v", err)
	}

	err = menu.RemoveMenu(id)
	if err != nil {
		t.Errorf("Error removing menu item: %v", err)
	}

	item := menuRepo.GetMenuItem(id)
	if item != nil {
		t.Errorf("Expected %v, got %v", nil, item)
	}
}
