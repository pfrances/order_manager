package model

import (
	"fmt"
	"order_manager/internal/id"
)

var (
	ErrPreparationNotFound      = fmt.Errorf("preparation not found")
	ErrPreparationNotPending    = fmt.Errorf("preparation is not pending")
	ErrPreparationNotInProgress = fmt.Errorf("preparation is not in progress")
	ErrPreparationNotReady      = fmt.Errorf("preparation is not ready")
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
