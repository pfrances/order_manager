package model

import (
	"fmt"
	"order_manager/internal/id"
)

var (
	ErrPreparationNotFound    = fmt.Errorf("preparation not found")
	ErrPreparationWrongStatus = fmt.Errorf("preparation has wrong status")
)

type PreparationStatus string

const (
	PreparationStatusInvalid    PreparationStatus = "invalid"
	PreparationStatusPending    PreparationStatus = "pending"
	PreparationStatusInProgress PreparationStatus = "in progress"
	PreparationStatusReady      PreparationStatus = "ready"
	PreparationStatusServed     PreparationStatus = "served"
	PreparationStatusAborted    PreparationStatus = "aborted"
)

type Preparation struct {
	ID         id.ID
	OrderID    id.ID
	MenuItemID id.ID
	Status     PreparationStatus
}
