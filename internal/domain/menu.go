package domain

import "order_manager/internal/id"

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

type MenuRepository interface {
	SaveItem(item MenuItem) error
	FindItem(id id.ID) (MenuItem, error)

	SaveCategory(category MenuCategory) error
	FindCategory(id id.ID) (MenuCategory, error)
}

type MenuService struct {
	repo MenuRepository
}

func NewMenuService(repo MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) CreateCategory(name string) (MenuCategory, error) {
	category := MenuCategory{
		ID:        id.NewID(),
		Name:      name,
		MenuItems: make([]MenuItem, 0),
	}

	err := s.repo.SaveCategory(category)
	if err != nil {
		return MenuCategory{}, err
	}

	return category, nil
}

func (s *MenuService) CreateMenuItem(name string, price int) (MenuItem, error) {
	item := MenuItem{
		ID:    id.NewID(),
		Name:  name,
		Price: price,
	}

	err := s.repo.SaveItem(item)
	if err != nil {
		return MenuItem{}, err
	}

	return item, nil
}

func (s *MenuService) AddItemToCategory(categoryID id.ID, itemID id.ID) error {
	category, err := s.repo.FindCategory(categoryID)
	if err != nil {
		return err
	}

	item, err := s.repo.FindItem(itemID)
	if err != nil {
		return err
	}

	category.MenuItems = append(category.MenuItems, item)
	return s.repo.SaveCategory(category)
}
