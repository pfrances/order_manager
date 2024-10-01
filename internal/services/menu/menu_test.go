package menu_test

import (
	"order_manager/internal/repositories/memory"
	"order_manager/internal/services/menu"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddItem(t *testing.T) {
	asserts := assert.New(t)
	repo := memory.NewMenuRepository()
	service := menu.NewService(repo)

	id, err := service.CreateMenu("Cheeseburger", 10)
	asserts.Nil(err)

	item := repo.GetMenuItem(id)
	asserts.NotNil(item)
}

func TestRemoveItem(t *testing.T) {
	asserts := assert.New(t)
	repo := memory.NewMenuRepository()
	service := menu.NewService(repo)

	id, err := service.CreateMenu("Cheeseburger", 10)
	asserts.Nil(err)

	err = service.RemoveMenu(id)
	asserts.Nil(err)

	item := repo.GetMenuItem(id)
	asserts.Nil(item)
}
