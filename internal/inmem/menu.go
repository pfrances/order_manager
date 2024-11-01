package inmem

import (
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"sync"
)

type Menu struct {
	categories map[id.ID]domain.MenuCategory
	items      map[id.ID]domain.MenuItem
	mu         sync.Mutex
}

func NewMenu() *Menu {
	return &Menu{
		categories: make(map[id.ID]domain.MenuCategory),
		items:      make(map[id.ID]domain.MenuItem),
	}
}

func (m *Menu) SaveItem(item domain.MenuItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[item.ID] = item
	return nil
}

func (m *Menu) FindItem(id id.ID) (domain.MenuItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.items[id]
	if !ok {
		return domain.MenuItem{}, domain.Errorf(domain.ENOTFOUND, "menu item with id %s not found", id)
	}
	return item, nil
}

func (m *Menu) SaveCategory(category domain.MenuCategory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.categories[category.ID] = category
	return nil
}

func (m *Menu) FindCategory(id id.ID) (domain.MenuCategory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	category, ok := m.categories[id]
	if !ok {
		return domain.MenuCategory{}, domain.Errorf(domain.ENOTFOUND, "menu category with id %s not found", id)
	}
	return category, nil
}
