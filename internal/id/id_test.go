package id_test

import (
	"order_manager/internal/id"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIDCollision(t *testing.T) {
	asserts := assert.New(t)
	firstID := id.NewID()
	secondID := id.NewID()

	asserts.NotEqual(firstID, secondID)
}
