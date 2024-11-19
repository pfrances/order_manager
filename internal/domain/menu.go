package domain

import (
	"context"
	"order_manager/internal/id"
	"slices"
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
	return i.ID != id.NilID() && i.Name != "" && i.Price >= 0
}

func (c MenuCategory) IsValid() bool {
	isValid := c.ID != id.NilID() && c.Name != "" && c.MenuItems != nil

	for _, item := range c.MenuItems {
		if !item.IsValid() {
			return false
		}
	}

	return isValid
}

type MenuRepository interface {
	SaveItem(ctx context.Context, item MenuItem) error
	FindItem(ctx context.Context, id id.ID) (MenuItem, error)
	FindItems(ctx context.Context, ids []id.ID) ([]MenuItem, error)
	FindAllItems(ctx context.Context) ([]MenuItem, error)

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

func (s *MenuService) FindAllMenuItems(ctx context.Context) ([]MenuItem, error) {
	return s.repo.FindAllItems(ctx)
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

	if !item.IsValid() {
		return MenuItem{}, Errorf(EINVALID, "invalid item")
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

	if slices.ContainsFunc(category.MenuItems, func(item MenuItem) bool { return item.ID == itemID }) {
		return Errorf(EINVALID, "item already exists in category")
	}

	item, err := s.repo.FindItem(ctx, itemID)
	if err != nil {
		return err
	}

	category.MenuItems = append(category.MenuItems, item)
	return s.repo.SaveCategory(ctx, category)
}
