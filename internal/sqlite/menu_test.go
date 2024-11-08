package sqlite_test

import (
	"context"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"order_manager/internal/sqlite"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GenerateDummyItem() domain.MenuItem {
	return domain.MenuItem{
		ID:    id.New(),
		Name:  "item",
		Price: 100,
	}
}

func GenerateDummyCategory() domain.MenuCategory {
	return domain.MenuCategory{
		ID:        id.New(),
		Name:      "category",
		MenuItems: make([]domain.MenuItem, 0),
	}
}

func TestSaveAndRetrieveItemByID(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	item := GenerateDummyItem()
	menuRepo := sqlite.NewMenu(db)

	err := menuRepo.SaveItem(context.Background(), item)
	require.NoErrorf(t, err, "failed to save item: %v", err)

	gotItem, err := menuRepo.FindItem(context.Background(), item.ID)
	require.NoErrorf(t, err, "failed to retrieve item: %v", err)

	assert.Equal(t, item, gotItem)
}

func TestSaveAndRetrieveItemsByID(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	item1 := GenerateDummyItem()
	item2 := GenerateDummyItem()
	item3 := GenerateDummyItem()
	menuRepo := sqlite.NewMenu(db)

	err := menuRepo.SaveItem(context.Background(), item1)
	require.NoErrorf(t, err, "failed to save item: %v", err)

	err = menuRepo.SaveItem(context.Background(), item2)
	require.NoErrorf(t, err, "failed to save item: %v", err)

	err = menuRepo.SaveItem(context.Background(), item3)
	require.NoErrorf(t, err, "failed to save item: %v", err)

	gotItems, err := menuRepo.FindItems(context.Background(), []id.ID{item1.ID, item2.ID, item3.ID})
	require.NoErrorf(t, err, "failed to retrieve items: %v", err)

	assert.ElementsMatch(t, []domain.MenuItem{item1, item2, item3}, gotItems)
}

func TestSaveAndRetrieveItemsByEmptyID(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	item1 := GenerateDummyItem()
	item2 := GenerateDummyItem()
	item3 := GenerateDummyItem()
	menuRepo := sqlite.NewMenu(db)

	err := menuRepo.SaveItem(context.Background(), item1)
	require.NoErrorf(t, err, "failed to save item: %v", err)

	err = menuRepo.SaveItem(context.Background(), item2)
	require.NoErrorf(t, err, "failed to save item: %v", err)

	err = menuRepo.SaveItem(context.Background(), item3)
	require.NoErrorf(t, err, "failed to save item: %v", err)

	gotItems, err := menuRepo.FindItems(context.Background(), []id.ID{})
	require.NoErrorf(t, err, "failed to retrieve items: %v", err)

	assert.ElementsMatch(t, []domain.MenuItem{}, gotItems)
}

func TestNotFoundItem(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	menuRepo := sqlite.NewMenu(db)

	_, err := menuRepo.FindItem(context.Background(), id.New())
	assert.Equal(t, domain.ErrorCode(err), domain.ENOTFOUND)
}

func TestSaveMenuItemWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	menuRepo := sqlite.NewMenu(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := menuRepo.SaveItem(ctx, GenerateDummyItem())
	assert.Error(t, err)
}

func TestFindItemWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	menuRepo := sqlite.NewMenu(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := menuRepo.FindItem(ctx, id.New())
	assert.Error(t, err)
}

func TestFindItemsWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	menuRepo := sqlite.NewMenu(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := menuRepo.FindItems(ctx, []id.ID{id.New()})
	assert.Error(t, err)
}

func TestSaveItemsWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	menuRepo := sqlite.NewMenu(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := menuRepo.SaveItem(ctx, GenerateDummyItem())
	assert.Error(t, err)
}

func TestSaveAndRetrieveMenuCategory(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	category := GenerateDummyCategory()
	menuRepo := sqlite.NewMenu(db)

	err := menuRepo.SaveCategory(context.Background(), category)
	require.NoErrorf(t, err, "failed to save item: %v", err)

	gotCategory, err := menuRepo.FindCategory(context.Background(), category.ID)
	require.NoErrorf(t, err, "failed to retrieve item: %v", err)

	assert.Equal(t, category, gotCategory)
}

func TestNotFoundCategory(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	menuRepo := sqlite.NewMenu(db)

	_, err := menuRepo.FindCategory(context.Background(), id.New())
	assert.Equal(t, domain.ErrorCode(err), domain.ENOTFOUND)
}

func TestSaveCategoryWithItems(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	item1 := GenerateDummyItem()
	item2 := GenerateDummyItem()
	item3 := GenerateDummyItem()
	category := GenerateDummyCategory()
	category.MenuItems = []domain.MenuItem{item1, item2, item3}
	menuRepo := sqlite.NewMenu(db)
	err := menuRepo.SaveItems(context.Background(), category.MenuItems)
	require.NoErrorf(t, err, "failed to save items: %v", err)

	err = menuRepo.SaveCategory(context.Background(), category)
	require.NoErrorf(t, err, "failed to save category: %v", err)

	gotCategory, err := menuRepo.FindCategory(context.Background(), category.ID)
	require.NoErrorf(t, err, "failed to retrieve category: %v", err)

	assert.Equal(t, category, gotCategory)
}

func TestSaveCategoryWithUnsavedItems(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	item1 := GenerateDummyItem()
	item2 := GenerateDummyItem()
	item3 := GenerateDummyItem()
	category := GenerateDummyCategory()
	category.MenuItems = []domain.MenuItem{item1, item2, item3}
	menuRepo := sqlite.NewMenu(db)

	err := menuRepo.SaveCategory(context.Background(), category)
	assert.Error(t, err)
}

func TestSaveCategoryWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	menuRepo := sqlite.NewMenu(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := menuRepo.SaveCategory(ctx, GenerateDummyCategory())
	assert.Error(t, err)
}

func TestFindCategoryWithContextCancellation(t *testing.T) {
	db := MustOpenDB(t)
	defer MustCloseDB(t, db)

	menuRepo := sqlite.NewMenu(db)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := menuRepo.FindCategory(ctx, id.New())
	assert.Error(t, err)
}
