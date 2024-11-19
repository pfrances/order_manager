package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"strings"
)

type dbBillStatus string

const (
	dbBillStatusOpen   dbBillStatus = "pending"
	dbBillStatusClosed dbBillStatus = "partially paid"
	dbBillStatusPaid   dbBillStatus = "paid"
)

func (s dbBillStatus) IsValid() bool {
	return s == dbBillStatusOpen || s == dbBillStatusClosed || s == dbBillStatusPaid
}

type dbBill struct {
	id      id.ID        `db:"id"`
	tableID id.ID        `db:"table_id"`
	total   int          `db:"total"`
	paid    int          `db:"paid"`
	status  dbBillStatus `db:"status"`
}

func (b dbBill) IsValid() bool {
	return b.id != id.NilID() && b.tableID != id.NilID() && b.total >= 0 && b.paid >= 0 && b.status.IsValid()
}

type Bill struct {
	*DB
}

func NewBill(db *DB) *Bill {
	return &Bill{DB: db}
}

func (b *Bill) Save(ctx context.Context, bill domain.Bill) error {
	if !bill.IsValid() {
		return domain.Errorf(domain.EINVALID, "bill is invalid: %v", bill)
	}

	tx, err := b.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO bills (id, table_id, total, paid, status)
		VALUES (?, ?, ?, ?, ?)
	`, bill.ID, bill.TableID, bill.TotalAmount, bill.Paid, dbBillStatus(bill.Status))
	if err != nil {
		return fmt.Errorf("failed to insert bill: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO bill_menu_items (bill_id, menu_item_id)
		VALUES %s
	`, strings.Repeat(", (?, ?)", len(bill.Items))[2:])
	args := make([]interface{}, 0, len(bill.Items)*2)
	for _, item := range bill.Items {
		args = append(args, bill.ID, item.ID)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert bill menu items: %w", err)
	}

	return tx.Commit()
}

func (b *Bill) FindByID(ctx context.Context, id id.ID) (domain.Bill, error) {
	tx, err := b.BeginTx(ctx, nil)
	if err != nil {
		return domain.Bill{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var dbBill dbBill
	err = tx.QueryRowContext(ctx, `
		SELECT id, table_id, total, paid, status
		FROM bills
		WHERE id = ?
	`, id).Scan(&dbBill.id, &dbBill.tableID, &dbBill.total, &dbBill.paid, &dbBill.status)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Bill{}, domain.Errorf(domain.ENOTFOUND, "bill with id %s not found", id)
		}
		return domain.Bill{}, fmt.Errorf("failed to find bill: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT id, name, price
		FROM menu_items
		WHERE id IN (
			SELECT menu_item_id
			FROM preparations
			WHERE order_id IN (
				SELECT id
				FROM orders
				WHERE table_id = ?
			)
		)
	`, dbBill.tableID)
	if err != nil {
		return domain.Bill{}, fmt.Errorf("failed to query menu items: %w", err)
	}
	defer rows.Close()

	var items []domain.MenuItem
	for rows.Next() {
		var item domain.MenuItem
		if err = rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			return domain.Bill{}, fmt.Errorf("failed to scan menu item: %w", err)
		}
		items = append(items, item)
	}

	bill := domain.Bill{
		ID:          dbBill.id,
		TableID:     dbBill.tableID,
		Items:       items,
		Status:      domain.BillStatus(dbBill.status),
		TotalAmount: dbBill.total,
		Paid:        dbBill.paid,
	}

	return bill, tx.Commit()
}

func (b *Bill) FindByTableID(ctx context.Context, id id.ID) ([]domain.Bill, error) {
	tx, err := b.BeginTx(ctx, nil)
	if err != nil {
		return []domain.Bill{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
	SELECT id, table_id, total, paid, status
	FROM bills
	WHERE table_id = ?
	`, id)
	if err != nil {
		return []domain.Bill{}, fmt.Errorf("failed to query bills: %w", err)
	}

	defer rows.Close()

	var bills []domain.Bill
	for rows.Next() {
		var dbBill dbBill
		err = rows.Scan(&dbBill.id, &dbBill.tableID, &dbBill.total, &dbBill.paid, &dbBill.status)
		if err != nil {
			return []domain.Bill{}, fmt.Errorf("failed to find bill: %w", err)
		}

		rows, err = tx.QueryContext(ctx, `
			SELECT id, name, price
			FROM menu_items
			WHERE id IN (
				SELECT menu_item_id
				FROM bill_menu_items
				WHERE bill_id = ?
			)
			`, dbBill.id)
		if err != nil {
			return []domain.Bill{}, fmt.Errorf("failed to query menu items: %w", err)
		}
		defer rows.Close()

		var items []domain.MenuItem
		for rows.Next() {
			var item domain.MenuItem
			if err = rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
				return []domain.Bill{}, fmt.Errorf("failed to scan menu item: %w", err)
			}
			items = append(items, item)
		}

		bills = append(bills, domain.Bill{
			ID:          dbBill.id,
			TableID:     dbBill.tableID,
			Status:      domain.BillStatus(dbBill.status),
			TotalAmount: dbBill.total,
			Paid:        dbBill.paid,
			Items:       items,
		})
	}

	return bills, tx.Commit()
}
