package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type menuItem struct {
	ID    id.ID
	Name  string
	Price int
}

func newMenuItem(m model.MenuItem) menuItem {
	return menuItem{
		ID:    m.ID,
		Name:  m.Name,
		Price: m.Price,
	}
}

func (m menuItem) toModel() model.MenuItem {
	return model.MenuItem{
		ID:    m.ID,
		Name:  m.Name,
		Price: m.Price,
	}
}

type MenuCategory struct {
	ID          id.ID
	Name        string
	MenuItemIds []id.ID
}
