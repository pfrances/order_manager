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
