package memory

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type MenuRepository struct {
	menu     map[id.ID]menuItem
	category map[id.ID]menuCategory
}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{
		menu:     make(map[id.ID]menuItem),
		category: make(map[id.ID]menuCategory),
	}
}

/* Menu */

func (m *MenuRepository) CreateMenuItem(menu model.MenuItem) error {
	m.menu[menu.ID] = menuItemFromModel(menu)
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

	m.menu[id] = menuItemFromModel(item)
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

func (m *MenuRepository) CreateMenuCategory(category model.MenuCategory) error {
	if _, ok := m.category[category.ID]; ok {
		return fmt.Errorf("category already exists")
	}

	m.category[category.ID] = menuCategoryFromModel(category)
	return nil
}

func (m *MenuRepository) GetMenuCategory(id id.ID) *model.MenuCategory {
	for _, category := range m.category {
		if category.ID == id {
			modelCategory := category.toModel()
			return &modelCategory
		}
	}
	return nil
}

func (m *MenuRepository) RemoveCategory(id id.ID) error {
	if _, ok := m.category[id]; !ok {
		return fmt.Errorf("category not found")
	}

	delete(m.category, id)
	return nil
}

func (m *MenuRepository) UpdateMenuCategory(id id.ID, fn func(menu *model.MenuCategory) error) error {
	category := m.GetMenuCategory(id)
	if category == nil {
		return fmt.Errorf("category not found")
	}

	if err := fn(category); err != nil {
		return err
	}

	m.category[id] = menuCategoryFromModel(*category)
	return nil
}

func (m *MenuRepository) GetMenuItemsInCategory(categoryID id.ID) []model.MenuItem {
	category := m.GetMenuCategory(categoryID)
	if category == nil {
		return nil
	}

	var items []model.MenuItem
	for _, itemID := range category.MenuItemIds {
		item := m.GetMenuItem(itemID)
		if item != nil {
			items = append(items, *item)
		}
	}

	return items
}
