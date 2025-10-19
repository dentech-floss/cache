package cache

import (
	"context"
	"testing"
	"time"
)

func TestNoOpCache(t *testing.T) {
	cache := NewNoOp[TestUser]()

	ctx := context.Background()
	user := TestUser{ID: "123", Name: "John"}

	// Test Set - should always succeed
	err := cache.Set(ctx, "key1", user, time.Minute)
	if err != nil {
		t.Errorf("Set should not return error, got: %v", err)
	}

	// Test Get - should always return false
	_, found := cache.Get(ctx, "key1")
	if found {
		t.Error("Get should always return false for NoOp cache")
	}

	// Test Delete - should always succeed
	err = cache.Delete(ctx, "key1")
	if err != nil {
		t.Errorf("Delete should not return error, got: %v", err)
	}

	// Test Close - should always succeed
	err = cache.Close()
	if err != nil {
		t.Errorf("Close should not return error, got: %v", err)
	}
}
