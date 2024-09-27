package model

import (
	"order_manager/internal/id"
)

type PreparationStatus string

const (
	PreparationStatusPending    PreparationStatus = "pending"
	PreparationStatusInProgress PreparationStatus = "in progress"
	PreparationStatusReady      PreparationStatus = "ready"
	PreparationStatusDone       PreparationStatus = "done"
)

type Preparation struct {
	ID      id.ID
	OrderID id.ID
	Status  PreparationStatus
}
