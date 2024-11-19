package domain_test

import (
	"context"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"order_manager/internal/inmem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsMenuItemValid(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			item     domain.MenuItem
		}{
			{
				testName: "valid item",
				item:     domain.MenuItem{ID: id.New(), Name: "test", Price: 100},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				assert.Truef(t, tc.item.IsValid(), "should be valid item: %v", tc.item)
			})
		}
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			item     domain.MenuItem
		}{
			{
				testName: "invalid item ID",
				item:     domain.MenuItem{ID: id.NilID(), Name: "test", Price: 100},
			},
			{
				testName: "empty item name",
				item:     domain.MenuItem{ID: id.New(), Name: "", Price: 100},
			},
			{
				testName: "negative item price",
				item:     domain.MenuItem{ID: id.New(), Name: "test", Price: -100},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				assert.Falsef(t, tc.item.IsValid(), "should be invalid item: %v", tc.item)
			})
		}
	})
}

func TestIsMenuCategoryValid(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName  string
			category  domain.MenuCategory
			menuItems []domain.MenuItem
		}{
			{
				testName:  "valid category with no items",
				category:  domain.MenuCategory{ID: id.New(), Name: "test", MenuItems: make([]domain.MenuItem, 0)},
				menuItems: make([]domain.MenuItem, 0),
			},
			{
				testName: "valid category with items",
				category: domain.MenuCategory{
					ID:        id.New(),
					Name:      "test",
					MenuItems: []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
				},
				menuItems: []domain.MenuItem{{ID: id.New(), Name: "test", Price: 100}},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				assert.Truef(t, tc.category.IsValid(), "should be valid category: %v", tc.category)
			})
		}
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName  string
			category  domain.MenuCategory
			menuItems []domain.MenuItem
		}{
			{
				testName:  "invalid category ID",
				category:  domain.MenuCategory{ID: id.NilID(), Name: "test", MenuItems: make([]domain.MenuItem, 0)},
				menuItems: make([]domain.MenuItem, 0),
			},
			{
				testName:  "empty category name",
				category:  domain.MenuCategory{ID: id.New(), Name: "", MenuItems: make([]domain.MenuItem, 0)},
				menuItems: make([]domain.MenuItem, 0),
			},
			{
				testName: "category with nil items",
				category: domain.MenuCategory{ID: id.New(), Name: "test", MenuItems: nil},
			},
			{
				testName: "category with invalid item",
				category: domain.MenuCategory{
					ID:        id.New(),
					Name:      "test",
					MenuItems: []domain.MenuItem{{ID: id.NilID(), Name: "test", Price: 100}},
				},
				menuItems: []domain.MenuItem{{ID: id.NilID(), Name: "test", Price: 100}},
			},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				t.Parallel()
				assert.Falsef(t, tc.category.IsValid(), "should be invalid category: %v", tc.category)
			})
		}
	})
}

func TestCreateMenuItem(t *testing.T) {
	menuRepo := inmem.NewMenu()
	menuService := domain.NewMenuService(menuRepo)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName  string
			itemName  string
			itemPrice int
		}{
			{testName: "valid item", itemName: "Spaghetti", itemPrice: 100},
			{testName: "item with price 0", itemName: "Spaghetti", itemPrice: 0},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				item, err := menuService.CreateMenuItem(context.Background(), tc.itemName, tc.itemPrice)
				require.Nil(t, err, "item creation failed")
				item, err = menuRepo.FindItem(context.Background(), item.ID)
				require.Nil(t, err, "item not saved")
				assert.Equal(t, tc.itemName, item.Name, "item name not correctly saved")
				assert.Equal(t, tc.itemPrice, item.Price, "item price not correctly saved")
			})
		}
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName  string
			itemName  string
			itemPrice int
		}{
			{testName: "empty item name", itemName: "", itemPrice: 100},
			{testName: "negative item price", itemName: "Spaghetti", itemPrice: -100},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				_, err := menuService.CreateMenuItem(context.Background(), tc.itemName, tc.itemPrice)
				require.NotNil(t, err, "item creation should fail")
			})
		}

	})
}

func TestCreateMenuCategory(t *testing.T) {
	menuRepo := inmem.NewMenu()
	menuService := domain.NewMenuService(menuRepo)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			category string
		}{
			{testName: "valid category", category: "Pasta"},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				category, err := menuService.CreateCategory(context.Background(), tc.category)
				require.Nil(t, err, "category creation failed")
				category, err = menuRepo.FindCategory(context.Background(), category.ID)
				require.Nil(t, err, "category not saved")
				assert.Equal(t, tc.category, category.Name, "name not correctly saved")
			})
		}
	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			testName string
			category string
		}{
			{testName: "empty category name", category: ""},
		}

		for _, tc := range tt {
			t.Run(tc.testName, func(t *testing.T) {
				_, err := menuService.CreateCategory(context.Background(), tc.category)
				require.NotNil(t, err, "category creation should fail")
			})
		}
	})
}

func TestAddMenuItemToCategory(t *testing.T) {
	menuRepo := inmem.NewMenu()
	menuService := domain.NewMenuService(menuRepo)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		category := domain.MenuCategory{ID: id.New(),
			Name:      "Pasta",
			MenuItems: make([]domain.MenuItem, 0),
		}
		err := menuRepo.SaveCategory(context.Background(), category)
		require.Nil(t, err, "error creating category")
		item := domain.MenuItem{ID: id.New(), Name: "Spaghetti", Price: 100}
		err = menuRepo.SaveItem(context.Background(), item)
		require.Nil(t, err, "error creating menu item")

		err = menuService.AddItemToCategory(context.Background(), category.ID, item.ID)

		require.NoError(t, err, "adding menu item to category failed")
		updatedCategory, err := menuRepo.FindCategory(context.Background(), category.ID)
		require.NoError(t, err, "category not found")
		assert.Equal(t, item.ID, updatedCategory.MenuItems[0].ID, "menu item not added to category")

	})

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		t.Run("category not found", func(t *testing.T) {
			t.Parallel()

			item := domain.MenuItem{ID: id.New(), Name: "Spaghetti", Price: 100}
			err := menuRepo.SaveItem(context.Background(), item)
			require.Nil(t, err, "error creating menu item")

			err = menuService.AddItemToCategory(context.Background(), id.New(), item.ID)
			require.Error(t, err, "adding menu item to category should fail")
		})

		t.Run("item not found", func(t *testing.T) {
			t.Parallel()

			category := domain.MenuCategory{ID: id.New(),
				Name:      "Pasta",
				MenuItems: make([]domain.MenuItem, 0),
			}
			err := menuRepo.SaveCategory(context.Background(), category)
			require.Nil(t, err, "error creating category")

			err = menuService.AddItemToCategory(context.Background(), category.ID, id.New())
			require.Error(t, err, "adding menu item to category should fail")
		})

		t.Run("item already in category", func(t *testing.T) {
			t.Parallel()

			item := domain.MenuItem{ID: id.New(), Name: "Spaghetti", Price: 100}
			category := domain.MenuCategory{ID: id.New(),
				Name:      "Pasta",
				MenuItems: []domain.MenuItem{item},
			}
			err := menuRepo.SaveCategory(context.Background(), category)
			require.Nil(t, err, "error creating category")
			err = menuRepo.SaveItem(context.Background(), item)
			require.Nil(t, err, "error creating menu item")

			err = menuService.AddItemToCategory(context.Background(), category.ID, item.ID)

			require.Error(t, err, "adding menu item to category should fail")
		})
	})

}
