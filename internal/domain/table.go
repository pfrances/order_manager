package domain

import (
	"context"
	"order_manager/internal/id"
)

type TableStatus string

const (
	TableStatusOpened TableStatus = "opened"
	TableStatusClosed TableStatus = "closed"
)

type Table struct {
	ID     id.ID
	Orders []Order
	Status TableStatus
}

type OrderStatus string

var (
	OrderStatusTaken   OrderStatus = "taken"
	OrderStatusDone    OrderStatus = "done"
	OrderStatusAborted OrderStatus = "aborted"
)

type Order struct {
	ID           id.ID
	Status       OrderStatus
	Preparations []Preparation
}

type PreparationStatus string

const (
	PreparationStatusPending    PreparationStatus = "pending"
	PreparationStatusInProgress PreparationStatus = "in progress"
	PreparationStatusReady      PreparationStatus = "ready"
	PreparationStatusServed     PreparationStatus = "served"
	PreparationStatusAborted    PreparationStatus = "aborted"
)

type Preparation struct {
	ID       id.ID
	MenuItem MenuItem
	Status   PreparationStatus
}

type TableRepository interface {
	Save(ctx context.Context, table Table) error
	FindByID(ctx context.Context, id id.ID) (Table, error)
	FindByStatus(ctx context.Context, status TableStatus) ([]Table, error)
}

type TableService struct {
	repo TableRepository
}

// NewTableService creates a new table service.
// The service is responsible for handling table related operations:
// such as opening and closing tables, taking orders, and managing preparations.
func NewTableService(repo TableRepository) *TableService {
	return &TableService{repo: repo}
}

// FindTable returns a table by its ID.
// Possible errors:
// - ENOTFOUND if the table could not be found.
func (s *TableService) FindTable(ctx context.Context, tableID id.ID) (Table, error) {
	return s.repo.FindByID(ctx, tableID)
}

// FindOpenedTables returns all tables with an open status.
// Possible errors:
// - Any error returned by the repository when fetching the tables.
func (s *TableService) FindOpenedTables(ctx context.Context) ([]Table, error) {
	return s.repo.FindByStatus(ctx, TableStatusOpened)
}

// OpenTable creates a new table with an open status and saves it to the repository.
// Possible errors:
// - Any error returned by the repository when saving the table.
func (s *TableService) OpenTable(ctx context.Context) (Table, error) {
	table := Table{
		ID:     id.NewID(),
		Status: TableStatusOpened,
		Orders: make([]Order, 0),
	}

	err := s.repo.Save(ctx, table)
	if err != nil {
		return Table{}, err
	}

	return table, nil
}

// CloseTable closes a table by setting its status to closed.
// Possible errors:
// - ENOTFOUND if the table could not be found.
// - EINVALID if the table is already closed.
// - Any error returned by the repository when saving the table.
func (s *TableService) CloseTable(ctx context.Context, tableID id.ID) error {
	table, err := s.repo.FindByID(ctx, tableID)
	if err != nil {
		return err
	}

	if table.Status == TableStatusClosed {
		return Errorf(EINVALID, "table %s is already closed", tableID)
	}

	table.Status = TableStatusClosed

	err = s.repo.Save(ctx, table)
	if err != nil {
		return err
	}

	return nil
}

// TakeOrder creates a new order for a table with the given menu items.
// Possible errors:
// - ENOTFOUND if the table could not be found.
// - EINVALID if the table is not open.
// - EINVALID if any of the menu items are invalid or the slice is empty.
// - Any error returned by the repository when saving the table.
func (s *TableService) TakeOrder(ctx context.Context, tableID id.ID, menuItems []MenuItem) (Order, error) {

	if len(menuItems) == 0 {
		return Order{}, Errorf(EINVALID, "no menu items provided")
	}

	for _, item := range menuItems {
		if !item.IsValid() {
			return Order{}, Errorf(EINVALID, "invalid menu item %s", item.ID)
		}
	}

	table, err := s.repo.FindByID(ctx, tableID)
	if err != nil {
		return Order{}, err
	}

	if table.Status != TableStatusOpened {
		return Order{}, Errorf(EINVALID, "table %s is not open", tableID)
	}

	order := Order{
		ID:           id.NewID(),
		Status:       OrderStatusTaken,
		Preparations: make([]Preparation, 0, len(menuItems)),
	}

	for _, item := range menuItems {
		prep := Preparation{
			ID:       id.NewID(),
			MenuItem: item,
			Status:   PreparationStatusPending,
		}
		order.Preparations = append(order.Preparations, prep)
	}

	table.Orders = append(table.Orders, order)

	err = s.repo.Save(ctx, table)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

// StartPreparation sets the status of a preparation to in progress.
//
// Possible errors:
// - ENOTFOUND if the table could not be found.
// - ENOTFOUND if the order could not be found.
// - ENOTFOUND if the preparation could not be found.
// - EINVALID if the table is not open.
// - EINVALID if the preparation is not pending.
// - Any error returned by the repository when saving the table.
func (s *TableService) StartPreparation(ctx context.Context, tableID, orderID, preparationID id.ID) error {
	table, err := s.repo.FindByID(ctx, tableID)
	if err != nil {
		return err
	}

	if table.Status != TableStatusOpened {
		return Errorf(EINVALID, "table %s is not open", tableID)
	}

	order, err := table.findOrder(orderID)
	if err != nil {
		return err
	}

	prep, err := order.findPreparation(preparationID)
	if err != nil {
		return Errorf(ENOTFOUND, "preparation %s not found in order %s", preparationID, orderID)
	}

	if prep.Status != PreparationStatusPending {
		return Errorf(EINVALID, "preparation %s is not pending, preparation status is %s", preparationID, prep.Status)
	}

	prep.Status = PreparationStatusInProgress
	order.updatePreparation(prep)
	table.updateOrder(order)

	err = s.repo.Save(ctx, table)
	if err != nil {
		return err
	}

	return nil
}

// FinishPreparation sets the status of a preparation to ready.
// Possible errors:
// - ENOTFOUND if the table could not be found.
// - ENOTFOUND if the order could not be found.
// - ENOTFOUND if the preparation could not be found.
// - EINVALID if the table is not open.
// - EINVALID if the preparation is not in progress.
// - Any error returned by the repository when saving the table.
func (s *TableService) FinishPreparation(ctx context.Context, tableID, orderID, preparationID id.ID) error {
	table, err := s.repo.FindByID(ctx, tableID)
	if err != nil {
		return err
	}

	if table.Status != TableStatusOpened {
		return Errorf(EINVALID, "table %s is not open", tableID)
	}

	order, err := table.findOrder(orderID)
	if err != nil {
		return err
	}

	prep, err := order.findPreparation(preparationID)
	if err != nil {
		return err
	}

	if prep.Status != PreparationStatusInProgress {
		return Errorf(EINVALID, "preparation %s is not in progress, preparation status is %s", preparationID, prep.Status)
	}

	prep.Status = PreparationStatusReady
	order.updatePreparation(prep)
	table.updateOrder(order)

	err = s.repo.Save(ctx, table)
	if err != nil {
		return err
	}

	return nil
}

// ServePreparation sets the status of a preparation to served.
// Possible errors:
// - ENOTFOUND if the table could not be found.
// - ENOTFOUND if the order could not be found.
// - ENOTFOUND if the preparation could not be found.
// - EINVALID if the table is not open.
// - EINVALID if the preparation is not ready.
// - Any error returned by the repository when saving the table.
func (s *TableService) ServePreparation(ctx context.Context, tableID id.ID, orderID id.ID, prepID id.ID) error {
	table, err := s.repo.FindByID(ctx, tableID)
	if err != nil {
		return err
	}

	if table.Status != TableStatusOpened {
		return Errorf(EINVALID, "table %s is not open", tableID)
	}

	order, err := table.findOrder(orderID)
	if err != nil {
		return err
	}

	prep, err := order.findPreparation(prepID)
	if err != nil {
		return err
	}

	if prep.Status != PreparationStatusReady {
		return Errorf(EINVALID, "preparation %s is not ready, preparation status is %s", prepID, prep.Status)
	}

	prep.Status = PreparationStatusServed
	order.updatePreparation(prep)

	var allServed = true
	for _, p := range order.Preparations {
		if p.Status != PreparationStatusServed {
			allServed = false
			break
		}
	}

	if allServed {
		order.Status = OrderStatusDone
	}

	table.updateOrder(order)

	err = s.repo.Save(ctx, table)
	if err != nil {
		return err
	}

	return nil
}

func (t *Table) findOrder(id id.ID) (Order, error) {
	for _, order := range t.Orders {
		if order.ID == id {
			return order, nil
		}
	}
	return Order{}, Errorf(ENOTFOUND, "order with id %s not found", id)
}

func (t *Table) updateOrder(order Order) {
	for i, o := range t.Orders {
		if o.ID == order.ID {
			t.Orders[i] = order
			return
		}
	}
}

func (o *Order) findPreparation(id id.ID) (Preparation, error) {
	for _, prep := range o.Preparations {
		if prep.ID == id {
			return prep, nil
		}
	}
	return Preparation{}, Errorf(ENOTFOUND, "preparation with id %s not found", id)
}

func (o *Order) updatePreparation(prep Preparation) {
	for i, p := range o.Preparations {
		if p.ID == prep.ID {
			o.Preparations[i] = prep
			return
		}
	}
}
