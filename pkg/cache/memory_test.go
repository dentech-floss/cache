package cache

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	cache := NewMemory[TestUser](nil)
	defer cache.Close()

	ctx := context.Background()
	user := TestUser{ID: "123", Name: "John"}

	// Test Set
	err := cache.Set(ctx, "key1", user, time.Minute)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Test Get
	retrieved, found := cache.Get(ctx, "key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if retrieved.ID != user.ID || retrieved.Name != user.Name {
		t.Errorf("Expected %+v, got %+v", user, retrieved)
	}

	// Test Get non-existent key
	_, found = cache.Get(ctx, "nonexistent")
	if found {
		t.Error("Expected not to find nonexistent key")
	}

	// Test Delete
	err = cache.Delete(ctx, "key1")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	// Verify deletion
	_, found = cache.Get(ctx, "key1")
	if found {
		t.Error("Expected key1 to be deleted")
	}
}

func TestMemoryCacheWithConfig(t *testing.T) {
	config := &MemoryConfig{
		SkipTTLExtensionOnHit: false, // Allow TTL extension on hit
	}
	cache := NewMemory[TestUser](config)
	defer cache.Close()

	ctx := context.Background()
	user := TestUser{ID: "123", Name: "John"}

	// Test Set
	err := cache.Set(ctx, "key1", user, 100*time.Millisecond)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Test Get before expiry
	retrieved, found := cache.Get(ctx, "key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if retrieved.ID != user.ID {
		t.Errorf("Expected ID %s, got %s", user.ID, retrieved.ID)
	}

	// Wait for expiry
	time.Sleep(150 * time.Millisecond)

	// Test Get after expiry
	_, found = cache.Get(ctx, "key1")
	if found {
		t.Error("Expected key1 to be expired")
	}
}

func TestMemoryCacheContextCancellation(t *testing.T) {
	cache := NewMemory[TestUser](nil)
	defer cache.Close()

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	user := TestUser{ID: "123", Name: "John"}

	// Test Set with cancelled context
	err := cache.Set(ctx, "key1", user, time.Minute)
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got: %v", err)
	}

	// Test Get with cancelled context
	_, found := cache.Get(ctx, "key1")
	if found {
		t.Error("Expected Get to return false with cancelled context")
	}

	// Test Delete with cancelled context
	err = cache.Delete(ctx, "key1")
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got: %v", err)
	}
}
