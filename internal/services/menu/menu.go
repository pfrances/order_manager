package menu

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type Service struct {
	repo repositories.MenuRepository
}

func NewService(repo repositories.MenuRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (m *Service) CreateMenu(name string, price int) (id.ID, error) {
	id := id.NewID()

	menu := model.MenuItem{
		ID:    id,
		Name:  name,
		Price: price,
	}
	return id, m.repo.CreateMenu(menu)
}

func (m *Service) GetMenuItem(id id.ID) *model.MenuItem {
	return m.repo.GetMenuItem(id)
}

func (m *Service) RemoveMenu(id id.ID) error {
	return m.repo.RemoveItem(id)
}

func (m *Service) CreateMenuCategory(name string) (id.ID, error) {
	id := id.NewID()

	category := model.MenuCategory{
		ID:   id,
		Name: name,
	}

	return id, m.repo.CreateMenuCategory(category)
}

func (m *Service) GetMenuCategory(id id.ID) *model.MenuCategory {
	return m.repo.GetMenuCategory(id)
}

func (m *Service) RemoveCategory(id id.ID) error {
	return m.repo.RemoveCategory(id)
}

func (m *Service) AddMenuItemToCategory(categoryID id.ID, itemID id.ID) error {
	return m.repo.UpdateMenuCategory(categoryID, func(category *model.MenuCategory) error {
		category.MenuItemIds = append(category.MenuItemIds, itemID)
		return nil
	})
}
