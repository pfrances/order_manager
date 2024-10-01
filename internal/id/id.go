package id

import "sync"

type ID int

var lastID = 0
var mutex = sync.Mutex{}

func NewID() ID {
	mutex.Lock()
	defer mutex.Unlock()

	newID := ID(lastID)
	lastID++

	return newID
}

func NilID() ID {
	return ID(-1)
}

func (id ID) IsNil() bool {
	return id == NilID()
}
