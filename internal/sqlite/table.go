package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"strings"
)

type dbTableStatus string

const (
	dbTableStatusOpened dbTableStatus = "opened"
	dbTableStatusClosed dbTableStatus = "closed"
)

type dbTable struct {
	id     id.ID         `db:"id"`
	status dbTableStatus `db:"status"`
}

type dbOrderStatus string

const (
	dbOrderStatusTaken   dbOrderStatus = "taken"
	dbOrderStatusDone    dbOrderStatus = "done"
	dbOrderStatusAborted dbOrderStatus = "aborted"
)

type dbOrder struct {
	id      id.ID         `db:"id"`
	tableID id.ID         `db:"table_id"`
	status  dbOrderStatus `db:"status"`
}

type dbPreparationStatus string

const (
	dbPreparationStatusPending    dbPreparationStatus = "pending"
	dbPreparationStatusInProgress dbPreparationStatus = "in progress"
	dbPreparationStatusReady      dbPreparationStatus = "ready"
	dbPreparationStatusServed     dbPreparationStatus = "served"
	dbPreparationStatusAborted    dbPreparationStatus = "aborted"
)

type dbPreparation struct {
	id         id.ID               `db:"id"`
	orderID    id.ID               `db:"order_id"`
	menuItemID id.ID               `db:"menu_item_id"`
	status     dbPreparationStatus `db:"status"`
}

type Table struct {
	*DB
}

func NewTable(db *DB) *Table {
	return &Table{DB: db}
}

func (t *Table) Save(ctx context.Context, table domain.Table) error {
	tx, err := t.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	dbTable, dbOrders, dbPreparations, err := toDBTable(table)
	if err != nil {
		return err
	}

	err = t.insertTable(ctx, tx, dbTable)
	if err != nil {
		return err
	}

	err = t.insertOrder(ctx, tx, dbOrders)
	if err != nil {
		return err
	}

	err = t.insertPreparation(ctx, tx, dbPreparations)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (t *Table) FindByID(ctx context.Context, id id.ID) (domain.Table, error) {
	tx, err := t.Begin()
	if err != nil {
		return domain.Table{}, err
	}
	defer tx.Rollback()

	var dbTable dbTable
	if err = tx.QueryRowContext(ctx, `
		SELECT id, status
		FROM tables
		WHERE id = ?
		`, id).Scan(&dbTable.id, &dbTable.status); err != nil {
		if err == sql.ErrNoRows {
			return domain.Table{}, domain.Errorf(domain.ENOTFOUND, "table %d not found", id)
		}
		return domain.Table{}, err
	}

	var dbOrders []dbOrder
	rows, err := tx.QueryContext(ctx, `
		SELECT id, table_id, status
		FROM orders
		WHERE table_id = ?
		`, id)
	if err != nil {
		return domain.Table{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var dbOrder dbOrder
		if err = rows.Scan(&dbOrder.id, &dbOrder.tableID, &dbOrder.status); err != nil {
			return domain.Table{}, fmt.Errorf("failed to scan order: %w", err)
		}
		dbOrders = append(dbOrders, dbOrder)
	}

	var dbPreparations []dbPreparation
	for _, o := range dbOrders {
		rows, err = tx.QueryContext(ctx, `
			SELECT id, order_id, menu_item_id, status
			FROM preparations
			WHERE order_id = ?
			`, o.id)
		if err != nil {
			return domain.Table{}, err
		}
		defer rows.Close()

		for rows.Next() {
			var dbPreparation dbPreparation
			if err = rows.Scan(&dbPreparation.id, &dbPreparation.orderID, &dbPreparation.menuItemID, &dbPreparation.status); err != nil {
				return domain.Table{}, err
			}
			dbPreparations = append(dbPreparations, dbPreparation)
		}
	}

	var dbItems []dbMenuItem
	for _, p := range dbPreparations {
		var dbItem dbMenuItem
		if err = tx.QueryRowContext(ctx, `
			SELECT id, name, price
			FROM menu_items
			WHERE id = ?
			`, p.menuItemID).Scan(&dbItem.id, &dbItem.name, &dbItem.price); err != nil {
			return domain.Table{}, err
		}
		dbItems = append(dbItems, dbItem)
	}

	table := toDomainTable(dbTable, dbOrders, dbPreparations, dbItems)

	return table, tx.Commit()
}

func (t *Table) FindByStatus(ctx context.Context, status domain.TableStatus) ([]domain.Table, error) {
	tx, err := t.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT id, status
		FROM tables
		WHERE status = ?
		`, dbTableStatus(status))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make([]domain.Table, 0)
	for rows.Next() {
		var dbTable dbTable
		if err = rows.Scan(&dbTable.id, &dbTable.status); err != nil {
			return nil, fmt.Errorf("failed to scan table: %w", err)
		}

		var dbOrders []dbOrder
		rows, err := tx.QueryContext(ctx, `
			SELECT id, table_id, status
			FROM orders
			WHERE table_id = ?
			`, dbTable.id)
		if err != nil {
			return nil, fmt.Errorf("failed to query orders: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var dbOrder dbOrder
			if err = rows.Scan(&dbOrder.id, &dbOrder.tableID, &dbOrder.status); err != nil {
				return nil, fmt.Errorf("failed to scan order: %w", err)
			}
			dbOrders = append(dbOrders, dbOrder)
		}

		var dbPreparations []dbPreparation
		for _, o := range dbOrders {
			rows, err = tx.QueryContext(ctx, `
				SELECT id, order_id, menu_item_id, status
				FROM preparations
				WHERE order_id = ?
				`, o.id)
			if err != nil {
				return nil, fmt.Errorf("failed to query preparations: %w", err)
			}
			defer rows.Close()

			for rows.Next() {
				var dbPreparation dbPreparation
				if err = rows.Scan(&dbPreparation.id, &dbPreparation.orderID, &dbPreparation.menuItemID, &dbPreparation.status); err != nil {
					return nil, fmt.Errorf("failed to scan preparation: %w", err)
				}
				dbPreparations = append(dbPreparations, dbPreparation)
			}
		}

		var dbItems []dbMenuItem
		for _, p := range dbPreparations {
			var dbItem dbMenuItem
			if err = tx.QueryRowContext(ctx, `
				SELECT id, name, price
				FROM menu_items
				WHERE id = ?
				`, p.menuItemID).Scan(&dbItem.id, &dbItem.name, &dbItem.price); err != nil {
				return nil, fmt.Errorf("failed to scan menu item: %w", err)
			}
			dbItems = append(dbItems, dbItem)
		}

		table := toDomainTable(dbTable, dbOrders, dbPreparations, dbItems)
		tables = append(tables, table)
	}

	return tables, tx.Commit()
}

func (t *Table) insertTable(ctx context.Context, tx *sql.Tx, table dbTable) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO tables (id, status)
		VALUES (?, ?)
		`, table.id, table.status)
	if err != nil {
		return fmt.Errorf("failed to insert table: %w", err)
	}

	return nil
}

func (t *Table) insertOrder(ctx context.Context, tx *sql.Tx, order []dbOrder) error {
	if len(order) == 0 {
		return nil
	}

	orderQuery := fmt.Sprintf(`
			INSERT INTO orders (id, table_id, status)
			VALUES %s
			`, strings.Repeat(", (?, ?, ?)", len(order))[2:])
	args := make([]interface{}, 0, len(order)*3)
	for _, o := range order {
		args = append(args, o.id, o.tableID, o.status)
	}

	_, err := tx.ExecContext(ctx, orderQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to insert orders: %w", err)
	}

	return nil
}

func (t *Table) insertPreparation(ctx context.Context, tx *sql.Tx, preparation []dbPreparation) error {
	if len(preparation) == 0 {
		return nil
	}

	preparationQuery := fmt.Sprintf(`
		INSERT INTO preparations (id, order_id, menu_item_id, status)
		VALUES %s
		`, strings.Repeat(", (?, ?, ?, ?)", len(preparation))[2:])
	args := make([]interface{}, 0, len(preparation)*4)
	for _, p := range preparation {
		args = append(args, p.id, p.orderID, p.menuItemID, p.status)
	}

	_, err := tx.ExecContext(ctx, preparationQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to insert preparations: %w", err)
	}

	return nil
}

func toDBTable(table domain.Table) (dbTable, []dbOrder, []dbPreparation, error) {
	dbTable := dbTable{
		id:     table.ID,
		status: dbTableStatus(table.Status),
	}

	dbOrders := make([]dbOrder, 0, len(table.Orders))
	dbPreparations := make([]dbPreparation, 0, len(table.Orders))
	for _, o := range table.Orders {
		dbOrder := dbOrder{
			id:      o.ID,
			tableID: table.ID,
			status:  dbOrderStatus(o.Status),
		}
		dbOrders = append(dbOrders, dbOrder)

		for _, p := range o.Preparations {
			dbPreparation := dbPreparation{
				id:         p.ID,
				orderID:    o.ID,
				menuItemID: p.MenuItem.ID,
				status:     dbPreparationStatus(p.Status),
			}
			dbPreparations = append(dbPreparations, dbPreparation)
		}
	}

	return dbTable, dbOrders, dbPreparations, nil
}

func toDomainTable(dbTable dbTable, dbOrders []dbOrder, dbPreparations []dbPreparation, dbItems []dbMenuItem) domain.Table {
	table := domain.Table{
		ID:     dbTable.id,
		Status: domain.TableStatus(dbTable.status),
		Orders: make([]domain.Order, 0, len(dbOrders)),
	}

	for _, o := range dbOrders {
		order := domain.Order{
			ID:           o.id,
			Status:       domain.OrderStatus(o.status),
			Preparations: make([]domain.Preparation, 0, len(dbPreparations)),
		}
		for _, p := range dbPreparations {
			if p.orderID != o.id {
				continue
			}

			var item domain.MenuItem
			for _, i := range dbItems {
				if i.id == p.menuItemID {
					item = domain.MenuItem{
						ID:    i.id,
						Name:  i.name,
						Price: i.price,
					}
					break
				}
			}
			preparation := domain.Preparation{
				ID:       p.id,
				MenuItem: item,
				Status:   domain.PreparationStatus(p.status),
			}
			order.Preparations = append(order.Preparations, preparation)
		}
		table.Orders = append(table.Orders, order)
	}

	return table
}
