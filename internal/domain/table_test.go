package domain_test

import (
	"order_manager/internal/domain"
	"order_manager/internal/id"
	"order_manager/internal/inmem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTable(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)

	table, err := tableService.CreateTable()

	require.NoError(t, err, "table creation failed")
	assert.NotEqual(t, table.ID, id.NilID(), "generated table ID is nil")
}

func TestTakeOrderSuccess(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	table := domain.Table{ID: id.NewID(), Orders: make([]domain.Order, 0)}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")
	menuItem := domain.MenuItem{ID: id.NewID(), Name: "test", Price: 100}

	order, err := tableService.TakeOrder(table.ID, []domain.MenuItem{menuItem})

	require.NoError(t, err, "take order failed")
	assert.NotEqual(t, id.NilID(), order.ID, "generated order ID is nil")
	assert.Equal(t, domain.OrderStatusTaken, order.Status, "invalid order status")
	assert.Equal(t, menuItem, order.Preparations[0].MenuItem, "preparation item mismatch")

	table, err = tableRepo.Find(table.ID)
	require.NoError(t, err, "table not correctly saved")
	assert.Equal(t, order, table.Orders[0], "order not correctly saved")
	assert.Equal(t, menuItem, table.Orders[0].Preparations[0].MenuItem, "preparation item not correctly saved")
}

func TestTakeOrderTableNotFound(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)

	_, err := tableService.TakeOrder(id.NewID(), []domain.MenuItem{{ID: id.NewID(), Name: "test", Price: 100}})

	require.Error(t, err, "take order should fail")
}

func TestStartPreparationSuccess(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	preparation := domain.Preparation{ID: id.NewID(), MenuItem: domain.MenuItem{ID: id.NewID(), Name: "test", Price: 100}, Status: domain.PreparationStatusPending}
	order := domain.Order{ID: id.NewID(), Status: domain.OrderStatusTaken, Preparations: []domain.Preparation{preparation}}
	table := domain.Table{ID: id.NewID(), Orders: []domain.Order{order}}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.StartPreparation(table.ID, order.ID, preparation.ID)

	require.NoError(t, err, "start preparation failed")
	updatedTable, err := tableRepo.Find(table.ID)
	require.NoError(t, err, "table not correctly saved")
	updatedOrder := updatedTable.Orders[0]
	updatedPreparation := updatedOrder.Preparations[0]
	assert.Equal(t, domain.PreparationStatusInProgress, updatedPreparation.Status, "preparation status not updated")
}

func TestStartPreparationOrderNotFound(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	table := domain.Table{ID: id.NewID(), Orders: make([]domain.Order, 0)}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.StartPreparation(table.ID, id.NewID(), id.NewID())

	require.Error(t, err, "start preparation should fail")
}

func TestStartPreparationPreparationNotFound(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	order := domain.Order{ID: id.NewID(), Status: domain.OrderStatusTaken, Preparations: make([]domain.Preparation, 0)}
	table := domain.Table{ID: id.NewID(), Orders: []domain.Order{order}}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.StartPreparation(table.ID, order.ID, id.NewID())

	require.Error(t, err, "start preparation should fail")
}

func TestStartPreparationPreparationNotPending(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	preparation := domain.Preparation{ID: id.NewID(), MenuItem: domain.MenuItem{ID: id.NewID(), Name: "test", Price: 100}, Status: domain.PreparationStatusInProgress}
	order := domain.Order{ID: id.NewID(), Status: domain.OrderStatusTaken, Preparations: []domain.Preparation{preparation}}
	table := domain.Table{ID: id.NewID(), Orders: []domain.Order{order}}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.StartPreparation(table.ID, order.ID, preparation.ID)

	require.Error(t, err, "start preparation should fail")
}

func TestFinishPreparationSuccess(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	preparation := domain.Preparation{ID: id.NewID(), MenuItem: domain.MenuItem{ID: id.NewID(), Name: "test", Price: 100}, Status: domain.PreparationStatusInProgress}
	order := domain.Order{ID: id.NewID(), Status: domain.OrderStatusTaken, Preparations: []domain.Preparation{preparation}}
	table := domain.Table{ID: id.NewID(), Orders: []domain.Order{order}}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.FinishPreparation(table.ID, order.ID, preparation.ID)

	require.NoError(t, err, "finish preparation failed")

	updatedTable, err := tableRepo.Find(table.ID)
	require.NoError(t, err, "table not correctly saved")
	updatedOrder := updatedTable.Orders[0]
	updatedPreparation := updatedOrder.Preparations[0]
	assert.Equal(t, domain.PreparationStatusReady, updatedPreparation.Status, "preparation status not updated")
}

func TestFinishPreparationOrderNotFound(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	table := domain.Table{ID: id.NewID(), Orders: make([]domain.Order, 0)}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.FinishPreparation(table.ID, id.NewID(), id.NewID())

	require.Error(t, err, "finish preparation should fail")
}

func TestFinishPreparationPreparationNotFound(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	order := domain.Order{ID: id.NewID(), Status: domain.OrderStatusTaken, Preparations: make([]domain.Preparation, 0)}
	table := domain.Table{ID: id.NewID(), Orders: []domain.Order{order}}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.FinishPreparation(table.ID, order.ID, id.NewID())

	require.Error(t, err, "finish preparation should fail")
}

func TestFinishPreparationPreparationNotInProgress(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	preparation := domain.Preparation{ID: id.NewID(), MenuItem: domain.MenuItem{ID: id.NewID(), Name: "test", Price: 100}, Status: domain.PreparationStatusReady}
	order := domain.Order{ID: id.NewID(), Status: domain.OrderStatusTaken, Preparations: []domain.Preparation{preparation}}
	table := domain.Table{ID: id.NewID(), Orders: []domain.Order{order}}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.FinishPreparation(table.ID, order.ID, preparation.ID)

	require.Error(t, err, "finish preparation should fail")
}

func TestServeLastPreparationSuccess(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	preparation := domain.Preparation{ID: id.NewID(), MenuItem: domain.MenuItem{ID: id.NewID(), Name: "test", Price: 100}, Status: domain.PreparationStatusReady}
	order := domain.Order{ID: id.NewID(), Status: domain.OrderStatusTaken, Preparations: []domain.Preparation{preparation}}
	table := domain.Table{ID: id.NewID(), Orders: []domain.Order{order}}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.ServePreparation(table.ID, order.ID, preparation.ID)

	require.NoError(t, err, "serve order failed")

	updatedTable, err := tableRepo.Find(table.ID)
	require.NoError(t, err, "table not correctly saved")
	updatedOrder := updatedTable.Orders[0]
	assert.Equal(t, domain.OrderStatusDone, updatedOrder.Status, "order status not updated")
}

func TestServeNotLastPreparationSuccess(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	preparation := domain.Preparation{ID: id.NewID(), MenuItem: domain.MenuItem{ID: id.NewID(), Name: "test", Price: 100}, Status: domain.PreparationStatusReady}
	preparation2 := domain.Preparation{ID: id.NewID(), MenuItem: domain.MenuItem{ID: id.NewID(), Name: "test2", Price: 200}, Status: domain.PreparationStatusReady}
	order := domain.Order{ID: id.NewID(), Status: domain.OrderStatusTaken, Preparations: []domain.Preparation{preparation, preparation2}}
	table := domain.Table{ID: id.NewID(), Orders: []domain.Order{order}}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.ServePreparation(table.ID, order.ID, preparation.ID)

	require.NoError(t, err, "serve order failed")
	updatedTable, err := tableRepo.Find(table.ID)
	require.NoError(t, err, "table not correctly saved")
	updatedOrder := updatedTable.Orders[0]
	assert.Equal(t, domain.OrderStatusTaken, updatedOrder.Status, "order status not updated")
}

func TestServePreparationOrderNotFound(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	table := domain.Table{ID: id.NewID(), Orders: make([]domain.Order, 0)}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.ServePreparation(table.ID, id.NewID(), id.NewID())

	require.Error(t, err, "serve order should fail")
}

func TestServePreparationPreparationNotFound(t *testing.T) {
	tableRepo := inmem.NewTable()
	tableService := domain.NewTableService(tableRepo)
	order := domain.Order{ID: id.NewID(), Status: domain.OrderStatusTaken, Preparations: make([]domain.Preparation, 0)}
	table := domain.Table{ID: id.NewID(), Orders: []domain.Order{order}}
	err := tableRepo.Save(table)
	require.NoError(t, err, "Initial setup failed")

	err = tableService.ServePreparation(table.ID, order.ID, id.NewID())

	require.Error(t, err, "serve order should fail")
}
