package cache

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestDistributedCacheWithTestcontainers(t *testing.T) {
	// Skip if Docker is not available
	if !isDockerAvailable() {
		t.Skip("Docker not available, skipping testcontainers test")
	}

	// Start Valkey container
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "valkey/valkey:7.2-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	valkeyContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start Valkey container: %v", err)
	}
	defer func(
		valkeyContainer testcontainers.Container,
		ctx context.Context,
		opts ...testcontainers.TerminateOption,
	) {
		_ = valkeyContainer.Terminate(ctx, opts...)
	}(valkeyContainer, ctx)

	// Get the connection details
	host, err := valkeyContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := valkeyContainer.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	addr := host + ":" + port.Port()

	// Test JSON serialization
	t.Run("JSON Serialization", func(t *testing.T) {
		config := &DistributedConfig{
			Addr:              addr,
			SerializationType: SerializationJSON,
		}

		cache, err := NewDistributedGeneric[TestUser](config)
		if err != nil {
			t.Fatalf("Failed to create distributed cache: %v", err)
		}
		defer func(cache Cache[TestUser]) {
			_ = cache.Close()
		}(cache)

		testCacheOperations(t, cache)
	})

	// Test Gob serialization
	t.Run("Gob Serialization", func(t *testing.T) {
		config := &DistributedConfig{
			Addr:              addr,
			SerializationType: SerializationGob,
		}

		cache, err := NewDistributedGeneric[TestUser](config)
		if err != nil {
			t.Fatalf("Failed to create distributed cache: %v", err)
		}
		defer func(cache Cache[TestUser]) {
			_ = cache.Close()
		}(cache)

		testCacheOperations(t, cache)
	})

	// Test with custom serializer
	t.Run("Custom Serializer", func(t *testing.T) {
		config := &DistributedConfig{
			Addr:       addr,
			Serializer: &JSONSerializer{},
		}

		cache, err := NewDistributedGeneric[TestUser](config)
		if err != nil {
			t.Fatalf("Failed to create distributed cache: %v", err)
		}
		defer func(cache Cache[TestUser]) {
			_ = cache.Close()
		}(cache)

		testCacheOperations(t, cache)
	})

	// Test health check
	t.Run("Health Check", func(t *testing.T) {
		config := &DistributedConfig{
			Addr:              addr,
			SerializationType: SerializationJSON,
		}

		cache, err := NewDistributedGeneric[TestUser](config)
		if err != nil {
			t.Fatalf("Failed to create distributed cache: %v", err)
		}
		defer func(cache Cache[TestUser]) {
			_ = cache.Close()
		}(cache)

		// Test Ping
		if healthChecker, ok := cache.(HealthChecker); ok {
			err := healthChecker.Ping(ctx)
			if err != nil {
				t.Errorf("Health check failed: %v", err)
			}
		} else {
			t.Error("Cache should implement HealthChecker interface")
		}
	})

	// Test TTL behavior
	t.Run("TTL Behavior", func(t *testing.T) {
		config := &DistributedConfig{
			Addr:              addr,
			SerializationType: SerializationJSON,
		}

		cache, err := NewDistributedGeneric[TestUser](config)
		if err != nil {
			t.Fatalf("Failed to create distributed cache: %v", err)
		}
		defer func(cache Cache[TestUser]) {
			_ = cache.Close()
		}(cache)

		user := TestUser{ID: "123", Name: "John"}

		// Set with short TTL
		err = cache.Set(ctx, "ttl-test", user, 100*time.Millisecond)
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}

		// Should be available immediately
		retrieved, found := cache.Get(ctx, "ttl-test")
		if !found {
			t.Error("Expected to find key immediately after set")
		}
		if retrieved.ID != user.ID {
			t.Errorf("Expected ID %s, got %s", user.ID, retrieved.ID)
		}

		// Wait for expiry
		time.Sleep(150 * time.Millisecond)

		// Should be expired
		_, found = cache.Get(ctx, "ttl-test")
		if found {
			t.Error("Expected key to be expired")
		}
	})

	// Test context cancellation
	t.Run("Context Cancellation", func(t *testing.T) {
		config := &DistributedConfig{
			Addr:              addr,
			SerializationType: SerializationJSON,
		}

		cache, err := NewDistributedGeneric[TestUser](config)
		if err != nil {
			t.Fatalf("Failed to create distributed cache: %v", err)
		}
		defer func(cache Cache[TestUser]) {
			_ = cache.Close()
		}(cache)

		// Test with cancelled context
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		user := TestUser{ID: "123", Name: "John"}

		// These should handle cancellation gracefully
		err = cache.Set(cancelledCtx, "key", user, time.Minute)
		if err != nil && err != context.Canceled {
			t.Errorf("Expected context.Canceled or nil, got: %v", err)
		}

		_, found := cache.Get(cancelledCtx, "key")
		if found {
			t.Error("Expected Get to return false with cancelled context")
		}

		err = cache.Delete(cancelledCtx, "key")
		if err != nil && err != context.Canceled {
			t.Errorf("Expected context.Canceled or nil, got: %v", err)
		}
	})
}

func TestDistributedCacheFactory(t *testing.T) {
	// Skip if Docker is not available
	if !isDockerAvailable() {
		t.Skip("Docker not available, skipping testcontainers test")
	}

	// Start Valkey container
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "valkey/valkey:7.2-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	valkeyContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start Valkey container: %v", err)
	}
	defer func(
		valkeyContainer testcontainers.Container,
		ctx context.Context,
		opts ...testcontainers.TerminateOption,
	) {
		_ = valkeyContainer.Terminate(ctx, opts...)
	}(valkeyContainer, ctx)

	// Get the connection details
	host, err := valkeyContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := valkeyContainer.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	addr := host + ":" + port.Port()

	// Test factory with distributed config
	config := &Config{
		Type: TypeDistributed,
		Distributed: &DistributedConfig{
			Addr:              addr,
			SerializationType: SerializationJSON,
		},
	}

	cache, err := New[TestUser](config)
	if err != nil {
		t.Fatalf("Factory failed to create distributed cache: %v", err)
	}
	defer func(cache Cache[TestUser]) {
		_ = cache.Close()
	}(cache)

	testCacheOperations(t, cache)
}

func TestDistributedCacheErrorHandling(t *testing.T) {
	// Test with invalid address
	config := &DistributedConfig{
		Addr:              "invalid:6379",
		SerializationType: SerializationJSON,
	}

	_, err := NewDistributedGeneric[TestUser](config)
	if err == nil {
		t.Error("Expected error for invalid address")
	}

	// Test with nil config
	_, err = NewDistributedGeneric[TestUser](nil)
	if err == nil {
		t.Error("Expected error for nil config")
	}

	// Test with invalid serialization type
	config = &DistributedConfig{
		Addr:              "localhost:6379",
		SerializationType: SerializationType("invalid"),
	}

	_, err = NewDistributedGeneric[TestUser](config)
	if err == nil {
		t.Error("Expected error for invalid serialization type")
	}
}

// Helper function to test basic cache operations
func testCacheOperations(
	t *testing.T,
	cache Cache[TestUser],
) {
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

// Helper function to check if Docker is available
func isDockerAvailable() bool {
	// Simple check - try to create a container request
	// This is a basic check, in practice you might want to ping Docker daemon
	return true // For now, assume Docker is available
}
