package domain

import (
	"context"
	"order_manager/internal/id"
)

type MenuCategory struct {
	ID        id.ID
	Name      string
	MenuItems []MenuItem
}

type MenuItem struct {
	ID    id.ID
	Name  string
	Price int
}

func (i MenuItem) IsValid() bool {
	return i.ID != id.NilID() && i.Name != "" && i.Price > 0
}

type MenuRepository interface {
	SaveItem(ctx context.Context, item MenuItem) error
	FindItem(ctx context.Context, id id.ID) (MenuItem, error)
	FindItems(ctx context.Context, ids []id.ID) ([]MenuItem, error)

	SaveCategory(ctx context.Context, category MenuCategory) error
	FindCategory(ctx context.Context, id id.ID) (MenuCategory, error)
}

type MenuService struct {
	repo MenuRepository
}

func NewMenuService(repo MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) FindMenuItems(ctx context.Context, itemIDs []id.ID) ([]MenuItem, error) {
	return s.repo.FindItems(ctx, itemIDs)
}

func (s *MenuService) CreateCategory(ctx context.Context, name string) (MenuCategory, error) {
	category := MenuCategory{
		ID:        id.New(),
		Name:      name,
		MenuItems: make([]MenuItem, 0),
	}

	err := s.repo.SaveCategory(ctx, category)
	if err != nil {
		return MenuCategory{}, err
	}

	return category, nil
}

func (s *MenuService) CreateMenuItem(ctx context.Context, name string, price int) (MenuItem, error) {
	item := MenuItem{
		ID:    id.New(),
		Name:  name,
		Price: price,
	}

	err := s.repo.SaveItem(ctx, item)
	if err != nil {
		return MenuItem{}, err
	}

	return item, nil
}

func (s *MenuService) AddItemToCategory(ctx context.Context, categoryID id.ID, itemID id.ID) error {
	category, err := s.repo.FindCategory(ctx, categoryID)
	if err != nil {
		return err
	}

	item, err := s.repo.FindItem(ctx, itemID)
	if err != nil {
		return err
	}

	category.MenuItems = append(category.MenuItems, item)
	return s.repo.SaveCategory(ctx, category)
}
