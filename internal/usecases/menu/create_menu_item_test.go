package menu_test

import (
	"order_manager/internal/repositories/memory"
	"order_manager/internal/usecases/menu"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMenuItemSuccess(t *testing.T) {
	menuRepo := memory.NewMenuRepository()
	usecase := menu.NewCreateMenuItem(menuRepo)

	id, err := usecase.Execute("Cheeseburger", 10)
	require.NoError(t, err)

	item := menuRepo.GetMenuItem(id)
	assert.Equal(t, "Cheeseburger", item.Name)
	assert.Equal(t, 10, item.Price)
}

func TestCreateMenuItemWithEmptyName(t *testing.T) {
	menuRepo := memory.NewMenuRepository()
	usecase := menu.NewCreateMenuItem(menuRepo)

	_, err := usecase.Execute("", 10)
	require.ErrorIs(t, err, menu.ErrMenuItemNameRequired)
}

func TestCreateMenuItemWithNegativePrice(t *testing.T) {
	menuRepo := memory.NewMenuRepository()
	usecase := menu.NewCreateMenuItem(menuRepo)

	_, err := usecase.Execute("Cheeseburger", -10)
	require.ErrorIs(t, err, menu.ErrMenuItemPriceInvalid)
}
