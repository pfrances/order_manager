package inmem

import (
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
		mu:     sync.Mutex{},
	}
}

func (t *Table) Save(table domain.Table) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.tables[table.ID] = table
	return nil
}

func (t *Table) Find(id id.ID) (domain.Table, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	table, ok := t.tables[id]
	if !ok {
		return domain.Table{}, domain.Errorf(domain.ENOTFOUND, "table with id %s not found", id)
	}
	return table, nil
}
