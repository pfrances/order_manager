package repositories

import (
	"context"
	"fmt"
	"order_manager/internal/id"
	"order_manager/internal/model"
)

var (
	ErrNotFound      = fmt.Errorf("not found")
	ErrAlreadyExists = fmt.Errorf("already exists")
)

type TransactionManager interface {
	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type KitchenRepository interface {
	CreatePreparation(preparation model.Preparation) error
	CreatePreparations(preparations []model.Preparation) error
	GetPreparation(id id.ID) *model.Preparation
	GetPreparationsByOrderID(orderID id.ID) []model.Preparation
	UpdatePreparation(id id.ID, fn func(preparation *model.Preparation) error) error
	RemovePreparation(id id.ID) error
}

// Order

type OrderRepository interface {
	CreateOrder(order model.Order) error
	GetOrder(id id.ID) *model.Order
	RemoveOrder(id id.ID) error
	UpdateOrder(id id.ID, fn func(order *model.Order) error) error
}

// Menu

type MenuRepository interface {
	CreateMenuItem(menu model.MenuItem) error
	GetMenuItem(id id.ID) *model.MenuItem
	UpdateMenuItem(id id.ID, fn func(menu *model.MenuItem) error) error
	RemoveItem(id id.ID) error
}

type MenuCategoryRepository interface {
	CreateMenuCategory(category model.MenuCategory) error
	GetMenuCategory(id id.ID) *model.MenuCategory
	UpdateMenuCategory(id id.ID, fn func(menu *model.MenuCategory) error) error
	RemoveCategory(id id.ID) error
}

// Table

type TableRepository interface {
	CreateTable(table model.Table) error
	GetTable(id id.ID) *model.Table
	UpdateTable(id id.ID, fn func(table *model.Table) error) error
	RemoveTable(id id.ID) error
}
