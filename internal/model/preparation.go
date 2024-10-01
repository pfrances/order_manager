package model

import (
	"fmt"
	"order_manager/internal/id"
)

var (
	ErrPreparationNotFound       = fmt.Errorf("preparation not found")
	ErrPreparationNotPending     = fmt.Errorf("preparation is not pending")
	ErrPreparationNotInProgress  = fmt.Errorf("preparation is not in progress")
	ErrPreparationAlreadyAborted = fmt.Errorf("preparation is already aborted")
)

type PreparationStatus string

const (
	PreparationStatusPending    PreparationStatus = "pending"
	PreparationStatusInProgress PreparationStatus = "in progress"
	PreparationStatusReady      PreparationStatus = "ready"
	PreparationStatusAborted    PreparationStatus = "aborted"
)

type Preparation struct {
	ID         id.ID
	OrderID    id.ID
	MenuItemID id.ID
	Status     PreparationStatus
}
