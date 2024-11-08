package id_test

import (
	"order_manager/internal/id"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIDCollision(t *testing.T) {
	firstID := id.New()
	secondID := id.New()

	assert.NotEqual(t, firstID, secondID)
}

func TestNilID(t *testing.T) {
	nilID := id.NilID()

	assert.True(t, nilID.IsNil())
}

func TestString(t *testing.T) {
	id := id.New()

	assert.NotEqual(t, id.String(), "")
}
