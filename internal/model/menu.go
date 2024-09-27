package model

import "order_manager/internal/id"

type MenuItem struct {
	ID    id.ID
	Name  string
	Price int
}

type MenuCategory struct {
	ID          id.ID
	Name        string
	MenuItemIds []id.ID
}
