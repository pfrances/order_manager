package menu

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type MenuService struct {
	repo Repository
}

type Repository interface {
	CreateMenu(menu *model.MenuItem) error
	GetMenuItem(id id.ID) *model.MenuItem
	UpdateMenuItem(id id.ID, fn func(menu *model.MenuItem) error) error
	RemoveItem(id id.ID) error

	CreateMenuCategory(category *model.MenuCategory) error
	GetMenuCategory(id id.ID) *model.MenuCategory
	UpdateMenuCategory(id id.ID, fn func(menu *model.MenuCategory) error) error
	RemoveCategory(id id.ID) error
}

func NewMenuService(repo Repository) *MenuService {
	return &MenuService{
		repo: repo,
	}
}

func (m *MenuService) CreateMenu(name string, price int) (id.ID, error) {
	id := id.NewID()

	menu := &model.MenuItem{
		ID:    id,
		Name:  name,
		Price: price,
	}
	return id, m.repo.CreateMenu(menu)
}

func (m *MenuService) GetMenuItem(id id.ID) *model.MenuItem {
	return m.repo.GetMenuItem(id)
}

func (m *MenuService) RemoveMenu(id id.ID) error {
	return m.repo.RemoveItem(id)
}

func (m *MenuService) CreateMenuCategory(name string) (id.ID, error) {
	id := id.NewID()

	category := &model.MenuCategory{
		ID:   id,
		Name: name,
	}

	return id, m.repo.CreateMenuCategory(category)
}

func (m *MenuService) GetMenuCategory(id id.ID) *model.MenuCategory {
	return m.repo.GetMenuCategory(id)
}

func (m *MenuService) RemoveCategory(id id.ID) error {
	return m.repo.RemoveCategory(id)
}

func (m *MenuService) AddMenuItemToCategory(categoryID id.ID, itemID id.ID) error {
	return m.repo.UpdateMenuCategory(categoryID, func(category *model.MenuCategory) error {
		category.MenuItemIds = append(category.MenuItemIds, itemID)
		return nil
	})
}
