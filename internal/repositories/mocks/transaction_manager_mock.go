package mocks

import "context"

type NewTransactionManagerMockOption func(*TransactionManagerMock)

type TransactionManagerMock struct {
	begin    func(ctx context.Context) (context.Context, error)
	commit   func(ctx context.Context) error
	rollback func(ctx context.Context) error
}

func WithBegin(begin func(ctx context.Context) (context.Context, error)) NewTransactionManagerMockOption {
	return func(m *TransactionManagerMock) {
		m.begin = begin
	}
}

func WithCommit(commit func(ctx context.Context) error) NewTransactionManagerMockOption {
	return func(m *TransactionManagerMock) {
		m.commit = commit
	}
}

func WithRollback(rollback func(ctx context.Context) error) NewTransactionManagerMockOption {
	return func(m *TransactionManagerMock) {
		m.rollback = rollback
	}
}

func NewTransactionManagerMock(options ...NewTransactionManagerMockOption) *TransactionManagerMock {
	manager := &TransactionManagerMock{}

	for _, option := range options {
		option(manager)
	}

	return manager
}

func (m *TransactionManagerMock) Begin(ctx context.Context) (context.Context, error) {
	return m.begin(ctx)
}

func (m *TransactionManagerMock) Commit(ctx context.Context) error {
	return m.commit(ctx)
}

func (m *TransactionManagerMock) Rollback(ctx context.Context) error {
	return m.rollback(ctx)
}
