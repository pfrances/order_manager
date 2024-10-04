package menu_test

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases/menu"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddMenuItemToCategorySuccess(t *testing.T) {
	menuRepo := memory.NewMenuRepository()
	usecase := menu.NewAddMenuItemToCategory(menuRepo)
	category := model.MenuCategory{
		ID:   id.NewID(),
		Name: "desserts",
	}
	menuItem := model.MenuItem{
		ID:    id.NewID(),
		Name:  "Cheesecake",
		Price: 10,
	}
	menuRepo.CreateMenuCategory(category)
	menuRepo.CreateMenuItem(menuItem)

	err := usecase.Execute(category.ID, menuItem.ID)
	require.NoError(t, err)

	updatedCategory := menuRepo.GetMenuCategory(category.ID)
	assert.Contains(t, updatedCategory.MenuItemIds, menuItem.ID)
}
