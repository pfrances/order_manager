package mocks

import (
	"order_manager/internal/id"
	"order_manager/internal/model"
)

type TableRepositoryMockOption func(*TableRepositoryMock)

type TableRepositoryMock struct {
	createTable func(table model.Table) error
	getTable    func(id id.ID) *model.Table
	updateTable func(id id.ID, fn func(table *model.Table) error) error
	removeTable func(id id.ID) error
}

func WithCreateTable(createTable func(table model.Table) error) TableRepositoryMockOption {
	return func(m *TableRepositoryMock) {
		m.createTable = createTable
	}
}

func WithGetTable(getTable func(id id.ID) *model.Table) TableRepositoryMockOption {
	return func(m *TableRepositoryMock) {
		m.getTable = getTable
	}
}

func WithUpdateTable(updateTable func(id id.ID, fn func(table *model.Table) error) error) TableRepositoryMockOption {
	return func(m *TableRepositoryMock) {
		m.updateTable = updateTable
	}
}

func WithRemoveTable(removeTable func(id id.ID) error) TableRepositoryMockOption {
	return func(m *TableRepositoryMock) {
		m.removeTable = removeTable
	}
}

func NewTableRepositoryMock(options ...TableRepositoryMockOption) *TableRepositoryMock {
	manager := &TableRepositoryMock{}

	for _, option := range options {
		option(manager)
	}

	return manager
}

func (m *TableRepositoryMock) CreateTable(table model.Table) error {
	return m.createTable(table)
}

func (m *TableRepositoryMock) GetTable(id id.ID) *model.Table {
	return m.getTable(id)
}

func (m *TableRepositoryMock) UpdateTable(id id.ID, fn func(table *model.Table) error) error {
	return m.updateTable(id, fn)
}

func (m *TableRepositoryMock) RemoveTable(id id.ID) error {
	return m.removeTable(id)
}
