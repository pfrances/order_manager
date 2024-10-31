package domain

import (
	"order_manager/internal/id"
)

type Table struct {
	ID     id.ID
	Orders []Order
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
	preparationStatusInvalid    PreparationStatus = "invalid"
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
	Save(table Table) error
	Find(id id.ID) (Table, error)
}

type TableService struct {
	repo TableRepository
}

func NewTableService(repo TableRepository) *TableService {
	return &TableService{repo: repo}
}

func (s *TableService) CreateTable() (Table, error) {
	table := Table{
		ID:     id.NewID(),
		Orders: make([]Order, 0),
	}

	err := s.repo.Save(table)
	if err != nil {
		return Table{}, err
	}

	return table, nil
}

func (s *TableService) TakeOrder(tableID id.ID, menuItem []MenuItem) (Order, error) {
	table, err := s.repo.Find(tableID)
	if err != nil {
		return Order{}, err
	}

	order := Order{
		ID:           id.NewID(),
		Status:       OrderStatusTaken,
		Preparations: make([]Preparation, 0, len(menuItem)),
	}

	for _, item := range menuItem {
		prep := Preparation{
			ID:       id.NewID(),
			MenuItem: item,
			Status:   PreparationStatusPending,
		}
		order.Preparations = append(order.Preparations, prep)
	}

	table.Orders = append(table.Orders, order)

	err = s.repo.Save(table)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}

func (s *TableService) StartPreparation(tableID, orderID, preparationID id.ID) error {
	table, err := s.repo.Find(tableID)
	if err != nil {
		return err
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

	err = s.repo.Save(table)
	if err != nil {
		return err
	}

	return nil
}

func (s *TableService) FinishPreparation(tableID, orderID, preparationID id.ID) error {
	table, err := s.repo.Find(tableID)
	if err != nil {
		return err
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

	err = s.repo.Save(table)
	if err != nil {
		return err
	}

	return nil
}

func (s *TableService) ServePreparation(tableID id.ID, orderID id.ID, prepID id.ID) error {
	table, err := s.repo.Find(tableID)
	if err != nil {
		return err
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

	err = s.repo.Save(table)
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
