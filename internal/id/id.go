package id

import (
	"strconv"
	"sync"
)

type ID int

var nextID = 1
var mutex sync.Mutex

func NewID() ID {
	mutex.Lock()
	defer mutex.Unlock()

	newID := ID(nextID)
	nextID++

	return newID
}

func NilID() ID {
	return ID(0)
}

func (id ID) IsNil() bool {
	return id == NilID()
}

func (id ID) String() string {
	return strconv.Itoa(int(id))
}
