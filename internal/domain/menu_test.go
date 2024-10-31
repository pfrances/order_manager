package domain_test

import (
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"order_manager/internal/inmem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddMenuItem(t *testing.T) {
	menuRepo := inmem.NewMenu()
	menuService := domain.NewMenuService(menuRepo)
	itemName := "Spaghetti"
	itemPrice := 100

	item, err := menuService.CreateMenuItem(itemName, itemPrice)

	require.Nil(t, err, "item creation failed")
	item, err = menuRepo.FindItem(item.ID)
	require.Nil(t, err, "item not saved")
	assert.Equal(t, itemName, item.Name, "item name not correctly saved")
	assert.Equal(t, itemPrice, item.Price, "item price not correctly saved")
}

func TestAddMenuCategory(t *testing.T) {
	menuRepo := inmem.NewMenu()
	menuService := domain.NewMenuService(menuRepo)
	categoryName := "Pasta"

	category, err := menuService.CreateCategory(categoryName)

	require.Nil(t, err, "category creation failed")
	category, err = menuRepo.FindCategory(category.ID)
	require.Nil(t, err, "category not saved")
	assert.Equal(t, categoryName, category.Name, "name not correctly saved")
}

func TestAddMenuItemToCategory(t *testing.T) {
	menuRepo := inmem.NewMenu()
	menuService := domain.NewMenuService(menuRepo)
	category := domain.MenuCategory{ID: id.NewID(), Name: "Pasta"}
	err := menuRepo.SaveCategory(category)
	require.Nil(t, err, "error creating category")
	item := domain.MenuItem{ID: id.NewID(), Name: "Spaghetti", Price: 100}
	err = menuRepo.SaveItem(item)
	require.Nil(t, err, "error creating menu item")

	err = menuService.AddItemToCategory(category.ID, item.ID)

	require.NoError(t, err, "adding menu item to category failed")
	updatedCategory, err := menuRepo.FindCategory(category.ID)
	require.NoError(t, err, "category not found")
	assert.Equal(t, item.ID, updatedCategory.MenuItems[0].ID, "menu item not added to category")
}

func TestAddMenuItemToCategoryNotFound(t *testing.T) {
	menuRepo := inmem.NewMenu()
	menuService := domain.NewMenuService(menuRepo)
	item := domain.MenuItem{ID: id.NewID(), Name: "Spaghetti", Price: 100}
	err := menuRepo.SaveItem(item)
	require.Nil(t, err, "error creating menu item")

	err = menuService.AddItemToCategory(id.NewID(), item.ID)

	require.Error(t, err, "adding menu item to category should fail")
}
