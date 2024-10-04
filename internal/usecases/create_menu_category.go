package usecases

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

var (
	ErrMenuCategoryNameRequired = fmt.Errorf("category name is required")
)

type CreateMenuCategory struct {
	menuRepository repositories.MenuRepository
}

func NewCreateMenuCategory(menuRepository repositories.MenuRepository) *CreateMenuCategory {
	return &CreateMenuCategory{menuRepository: menuRepository}
}

func (c *CreateMenuCategory) Execute(name string, menuItemID []id.ID) (id.ID, error) {
	if name == "" {
		return id.NilID(), ErrMenuCategoryNameRequired
	}

	category := model.MenuCategory{
		ID:          id.NewID(),
		Name:        name,
		MenuItemIds: menuItemID,
	}

	return category.ID, c.menuRepository.CreateMenuCategory(category)
}
