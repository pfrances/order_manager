package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type menuItem struct {
	id    id.ID
	name  string
	price int
}

func menuItemFromModel(m model.MenuItem) menuItem {
	return menuItem{
		id:    m.ID,
		name:  m.Name,
		price: m.Price,
	}
}

func (m menuItem) toModel() model.MenuItem {
	return model.MenuItem{
		ID:    m.id,
		Name:  m.name,
		Price: m.price,
	}
}

type menuCategory struct {
	ID          id.ID
	Name        string
	MenuItemIds []id.ID
}

func menuCategoryFromModel(m model.MenuCategory) menuCategory {
	return menuCategory{
		ID:          m.ID,
		Name:        m.Name,
		MenuItemIds: m.MenuItemIds,
	}
}

func (m menuCategory) toModel() model.MenuCategory {
	return model.MenuCategory{
		ID:          m.ID,
		Name:        m.Name,
		MenuItemIds: m.MenuItemIds,
	}
}
