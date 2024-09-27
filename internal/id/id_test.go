package id_test

import (
	"order_manager/internal/id"
	"testing"
)

func TestNewIDCollision(t *testing.T) {
	firstID := id.NewID()
	secondID := id.NewID()

	if firstID == secondID {
		t.Fatalf("Expected IDs to be different, got %v and %v", firstID, secondID)
	}
}
