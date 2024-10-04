package usecases_test

import (
	"order_manager/internal/id"
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMenuCategoryWithNoMenuItem(t *testing.T) {
	menuRepo := memory.NewMenuRepository()
	usecase := usecases.NewCreateMenuCategory(menuRepo)

	id, err := usecase.Execute("desserts", []id.ID{})
	require.NoError(t, err)

	category := menuRepo.GetMenuCategory(id)
	assert.Equal(t, "desserts", category.Name)
	assert.Empty(t, category.MenuItemIds)
}

func TestCreateMenuCategoryWithMenuItems(t *testing.T) {
	menuRepo := memory.NewMenuRepository()
	usecase := usecases.NewCreateMenuCategory(menuRepo)

	menuItemIDs := []id.ID{id.NewID(), id.NewID()}
	name := "desserts"
	id, err := usecase.Execute(name, menuItemIDs)
	require.NoError(t, err)

	category := menuRepo.GetMenuCategory(id)
	assert.Equal(t, name, category.Name)
	assert.Contains(t, category.MenuItemIds, menuItemIDs[0])
	assert.Contains(t, category.MenuItemIds, menuItemIDs[1])
}
