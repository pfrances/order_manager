package usecases

import (
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

var (
	ErrMenuItemNameRequired = fmt.Errorf("item name is required")
	ErrMenuItemPriceInvalid = fmt.Errorf("price must be greater than 0")
)

type CreateMenuItem struct {
	menuRepository repositories.MenuRepository
}

func NewCreateMenuItem(menuRepository repositories.MenuRepository) *CreateMenuItem {
	return &CreateMenuItem{menuRepository: menuRepository}
}

func (c *CreateMenuItem) Execute(name string, price int) (id.ID, error) {
	if name == "" {
		return id.NilID(), ErrMenuItemNameRequired
	}

	if price <= 0 {
		return id.NilID(), ErrMenuItemPriceInvalid
	}

	item := model.MenuItem{
		ID:    id.NewID(),
		Name:  name,
		Price: price,
	}

	return item.ID, c.menuRepository.CreateMenuItem(item)
}
