package inmem

import (
	"context"
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"sync"
)

type Table struct {
	tables map[id.ID]domain.Table
	mu     sync.Mutex
}

func NewTable() *Table {
	return &Table{
		tables: make(map[id.ID]domain.Table),
	}
}

func (t *Table) Save(ctx context.Context, table domain.Table) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.tables[table.ID] = table
	return nil
}

func (t *Table) FindByID(ctx context.Context, id id.ID) (domain.Table, error) {
	if ctx.Err() != nil {
		return domain.Table{}, ctx.Err()
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	table, ok := t.tables[id]
	if !ok {
		return domain.Table{}, domain.Errorf(domain.ENOTFOUND, "table with id %s not found", id)
	}
	return table, nil
}

func (t *Table) FindByStatus(ctx context.Context, status domain.TableStatus) ([]domain.Table, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	tables := make([]domain.Table, 0)
	for _, table := range t.tables {
		if table.Status == status {
			tables = append(tables, table)
		}
	}
	return tables, nil
}
