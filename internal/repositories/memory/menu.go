package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
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
	if _, ok := m.menu[menu.ID]; ok {
		return repositories.ErrAlreadyExists
	}

	m.menu[menu.ID] = menuItemFromModel(menu)
	return nil
}

func (m *MenuRepository) GetMenuItem(id id.ID) *model.MenuItem {
	for _, item := range m.menu {
		if item.id == id {
			modelItem := item.toModel()
			return &modelItem
		}
	}
	return nil
}

func (m *MenuRepository) UpdateMenuItem(id id.ID, fn func(menu *model.MenuItem) error) error {
	if _, ok := m.menu[id]; !ok {
		return repositories.ErrNotFound
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
		return repositories.ErrNotFound
	}

	delete(m.menu, id)
	return nil
}

/* MenuCategory */

func (m *MenuRepository) CreateMenuCategory(category model.MenuCategory) error {
	if _, ok := m.category[category.ID]; ok {
		return repositories.ErrAlreadyExists
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
		return repositories.ErrNotFound
	}

	delete(m.category, id)
	return nil
}

func (m *MenuRepository) UpdateMenuCategory(id id.ID, fn func(menu *model.MenuCategory) error) error {
	category := m.GetMenuCategory(id)
	if category == nil {
		return repositories.ErrNotFound
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
