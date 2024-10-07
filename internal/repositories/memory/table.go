package memory

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
	"order_manager/internal/repositories"
)

type TableRepository struct {
	table map[id.ID]table
}

func NewTableRepository() *TableRepository {
	return &TableRepository{table: make(map[id.ID]table)}
}

func (m *TableRepository) CreateTable(table model.Table) error {
	if _, ok := m.table[table.ID]; ok {
		return repositories.ErrAlreadyExists
	}

	m.table[table.ID] = tableFromModel(table)
	return nil
}

func (m *TableRepository) GetTable(id id.ID) *model.Table {
	table, ok := m.table[id]
	if !ok {
		return nil
	}

	modelTable := table.toModel()
	return &modelTable
}

func (m *TableRepository) RemoveTable(id id.ID) error {
	if _, ok := m.table[id]; !ok {
		return repositories.ErrNotFound
	}

	delete(m.table, id)
	return nil
}

func (m *TableRepository) UpdateTable(id id.ID, fn func(table *model.Table) error) error {
	table, ok := m.table[id]
	if !ok {
		return repositories.ErrNotFound
	}

	modelTable := table.toModel()
	err := fn(&modelTable)
	if err != nil {
		return err
	}

	m.table[id] = tableFromModel(modelTable)
	return nil
}
