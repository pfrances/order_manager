package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"strings"
)

type dbMenuItem struct {
	id    id.ID  `db:"id"`
	name  string `db:"name"`
	price int    `db:"price"`
}

func (i dbMenuItem) IsValid() bool {
	return i.id != id.NilID() && i.name != "" && i.price >= 0
}

type dbMenuItemCategory struct {
	id   id.ID  `db:"id"`
	name string `db:"name"`
}

func (c dbMenuItemCategory) IsValid() bool {
	return c.id != id.NilID() && c.name != ""
}

type dbMenuItemCategoryItem struct {
	categoryID id.ID `db:"category_id"`
	itemID     id.ID `db:"item_id"`
}

func (i dbMenuItemCategoryItem) IsValid() bool {
	return i.categoryID != id.NilID() && i.itemID != id.NilID()
}

type Menu struct {
	*DB
}

func NewMenu(db *DB) *Menu {
	return &Menu{DB: db}
}

func (m *Menu) SaveItem(ctx context.Context, item domain.MenuItem) error {
	if !item.IsValid() {
		return domain.Errorf(domain.EINVALID, "menu item is invalid: %v", item)
	}

	_, err := m.ExecContext(ctx, `
		INSERT INTO menu_items (id, name, price)
		VALUES (?, ?, ?)
	`, item.ID, item.Name, item.Price)
	if err != nil {
		return fmt.Errorf("failed to insert item: %w", err)
	}

	return nil
}

func (m *Menu) SaveItems(ctx context.Context, items []domain.MenuItem) error {
	if len(items) == 0 {
		return nil
	}

	for _, item := range items {
		if !item.IsValid() {
			return domain.Errorf(domain.EINVALID, "menu item is invalid: %v", item)
		}
	}

	tx, err := m.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	dbItems := make([]dbMenuItem, 0, len(items))
	for _, item := range items {
		dbItems = append(dbItems, dbMenuItem{
			id:    item.ID,
			name:  item.Name,
			price: item.Price,
		})
	}

	if err := m.insertItems(ctx, tx, dbItems); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *Menu) FindItem(ctx context.Context, id id.ID) (domain.MenuItem, error) {
	var item dbMenuItem
	err := m.QueryRowContext(ctx, `
		SELECT id, name, price
		FROM menu_items WHERE id = ?
		`, id).Scan(&item.id, &item.name, &item.price)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.MenuItem{}, domain.Errorf(domain.ENOTFOUND, "failed to find item with id %s", id)
		}
		return domain.MenuItem{}, fmt.Errorf("failed to find item: %w", err)
	}

	return domain.MenuItem{
		ID:    item.id,
		Name:  item.name,
		Price: item.price,
	}, nil
}

func (m *Menu) FindItems(ctx context.Context, ids []id.ID) ([]domain.MenuItem, error) {
	items := make([]domain.MenuItem, 0, len(ids))
	if len(ids) == 0 {
		return items, nil
	}

	query := fmt.Sprintf(`
		SELECT id, name, price 
		FROM menu_items
		WHERE id IN (%s)
		`, strings.Repeat(", ?", len(ids))[2:])
	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	rows, err := m.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item dbMenuItem
		if err := rows.Scan(&item.id, &item.name, &item.price); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}

		items = append(items, domain.MenuItem{
			ID:    item.id,
			Name:  item.name,
			Price: item.price,
		})
	}

	if len(items) != len(ids) {
		return nil, domain.Errorf(domain.ENOTFOUND, "failed to find all items")
	}

	return items, nil
}

func (m *Menu) FindAllItems(ctx context.Context) ([]domain.MenuItem, error) {
	rows, err := m.QueryContext(ctx, `
		SELECT id, name, price
		FROM menu_items
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}

	items := make([]domain.MenuItem, 0)
	for rows.Next() {
		var item dbMenuItem
		if err := rows.Scan(&item.id, &item.name, &item.price); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}

		items = append(items, domain.MenuItem{
			ID:    item.id,
			Name:  item.name,
			Price: item.price,
		})
	}

	return items, nil
}

func (m *Menu) SaveCategory(ctx context.Context, category domain.MenuCategory) error {
	if !category.IsValid() {
		return domain.Errorf(domain.EINVALID, "menu category is invalid: %v", category)
	}

	tx, err := m.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO menu_categories (id, name)
		VALUES (?, ?)
		`, category.ID, category.Name)
	if err != nil {
		return fmt.Errorf("failed to insert category: %w", err)
	}

	if len(category.MenuItems) == 0 {
		return tx.Commit()
	}

	menuItemCategoryQuery := fmt.Sprintf(`
		INSERT INTO menu_item_categories (category_id, item_id)
		VALUES %s
		`, strings.Repeat(", (?, ?)", len(category.MenuItems))[2:])
	args := make([]interface{}, 0, len(category.MenuItems)*2)
	for _, item := range category.MenuItems {
		args = append(args, category.ID, item.ID)
	}

	_, err = tx.ExecContext(ctx, menuItemCategoryQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to insert menu item categories: %w", err)
	}

	return tx.Commit()
}

func (m *Menu) FindCategory(ctx context.Context, id id.ID) (domain.MenuCategory, error) {
	tx, err := m.BeginTx(ctx, nil)
	if err != nil {
		return domain.MenuCategory{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var category dbMenuItemCategory
	err = tx.QueryRowContext(ctx, `
		SELECT id, name 
		FROM menu_categories 
		WHERE id = ?
		`, id).Scan(&category.id, &category.name)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.MenuCategory{}, domain.Errorf(domain.ENOTFOUND, "failed to find category with id %s", id)
		}
		return domain.MenuCategory{}, fmt.Errorf("failed to find category: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT id, name, price
		FROM menu_items 
		WHERE id 
		IN (
			SELECT item_id 
			FROM menu_item_categories 
			WHERE category_id = ?
		)
	`, id)
	if err != nil {
		return domain.MenuCategory{}, fmt.Errorf("failed to query menu item categories: %w", err)
	}
	defer rows.Close()

	items := make([]domain.MenuItem, 0)
	for rows.Next() {
		var item dbMenuItem
		if err := rows.Scan(&item.id, &item.name, &item.price); err != nil {
			return domain.MenuCategory{}, fmt.Errorf("failed to scan menu item: %w", err)
		}

		items = append(items, domain.MenuItem{
			ID:    item.id,
			Name:  item.name,
			Price: item.price,
		})
	}

	return domain.MenuCategory{
		ID:        category.id,
		Name:      category.name,
		MenuItems: items,
	}, tx.Commit()
}

func (m *Menu) insertItems(context context.Context, tx *sql.Tx, items []dbMenuItem) error {
	if len(items) == 0 {
		return nil
	}

	for _, i := range items {
		if !i.IsValid() {
			return domain.Errorf(domain.EINVALID, "menu item is invalid: %v", i)
		}
	}

	itemQuery := fmt.Sprintf(`
		INSERT INTO menu_items (id, name, price) 
		VALUES %s
		`, strings.Repeat(", (?, ?, ?)", len(items))[2:])
	args := make([]interface{}, 0, len(items)*3)
	for _, i := range items {
		args = append(args, i.id, i.name, i.price)
	}

	_, err := tx.ExecContext(context, itemQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to insert items: %w", err)
	}

	return nil
}
