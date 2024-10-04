package menu

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type AddMenuItemToCategory struct {
	menuRepository repositories.MenuRepository
}

func NewAddMenuItemToCategory(menuRepository repositories.MenuRepository) *AddMenuItemToCategory {
	return &AddMenuItemToCategory{menuRepository: menuRepository}
}

func (a *AddMenuItemToCategory) Execute(categoryID id.ID, menuItemID id.ID) error {
	return a.menuRepository.UpdateMenuCategory(categoryID, func(category *model.MenuCategory) error {
		category.MenuItemIds = append(category.MenuItemIds, menuItemID)
		return nil
	})
}
