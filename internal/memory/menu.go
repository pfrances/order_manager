package memory

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type MenuRepository struct {
	menu map[id.ID]menuItem
}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{}
}

/* Menu */

func (m *MenuRepository) CreateMenu(menu *model.MenuItem) error {
	m.menu[menu.ID] = newMenuItem(*menu)
	return nil
}

func (m *MenuRepository) GetMenuItem(id id.ID) *model.MenuItem {
	for _, item := range m.menu {
		if item.ID == id {
			modelItem := item.toModel()
			return &modelItem
		}
	}
	return nil
}

func (m *MenuRepository) UpdateMenuItem(id id.ID, fn func(menu *model.MenuItem) error) error {
	if _, ok := m.menu[id]; !ok {
		return fmt.Errorf("item not found")
	}

	item := m.menu[id].toModel()
	if err := fn(&item); err != nil {
		return err
	}

	m.menu[id] = newMenuItem(item)
	return nil
}

func (m *MenuRepository) RemoveItem(id id.ID) error {
	if _, ok := m.menu[id]; !ok {
		return fmt.Errorf("item not found")
	}

	delete(m.menu, id)
	return nil
}

/* MenuCategory */

func (m *MenuRepository) CreateMenuCategory(category *model.MenuCategory) error {
	// TODO
	return nil
}

func (m *MenuRepository) GetMenuCategory(id id.ID) *model.MenuCategory {
	// TODO
	return nil
}

func (m *MenuRepository) RemoveCategory(id id.ID) error {
	// TODO
	return nil
}

func (m *MenuRepository) UpdateMenuCategory(id id.ID, fn func(menu *model.MenuCategory) error) error {
	// TODO
	return nil
}

func (m *MenuRepository) GetMenuItemsInCategory(categoryID id.ID) []model.MenuItem {
	// TODO
	return nil
}
