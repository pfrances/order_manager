package menu

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
	menuCategoryRepository repositories.MenuCategoryRepository
}

func NewCreateMenuCategory(menuCategoryRepository repositories.MenuCategoryRepository) *CreateMenuCategory {
	return &CreateMenuCategory{menuCategoryRepository: menuCategoryRepository}
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

	return category.ID, c.menuCategoryRepository.CreateMenuCategory(category)
}
