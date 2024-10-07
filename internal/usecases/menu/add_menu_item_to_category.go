package menu

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type AddMenuItemToCategory struct {
	menuCategoryRepository repositories.MenuCategoryRepository
}

func NewAddMenuItemToCategory(menuCategoryRepository repositories.MenuCategoryRepository) *AddMenuItemToCategory {
	return &AddMenuItemToCategory{menuCategoryRepository: menuCategoryRepository}
}

func (a *AddMenuItemToCategory) Execute(categoryID id.ID, menuItemID id.ID) error {
	return a.menuCategoryRepository.UpdateMenuCategory(categoryID, func(category *model.MenuCategory) error {
		category.MenuItemIds = append(category.MenuItemIds, menuItemID)
		return nil
	})
}
