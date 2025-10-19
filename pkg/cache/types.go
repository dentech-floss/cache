package cache

import (
	"context"
	"time"
)

// Cache provides a generic key-value caching interface with type safety.
// Implementations can be in-memory, distributed, or no-op.
type Cache[T any] interface {
	// Get retrieves a value from the cache by key.
	// Returns the value and true if found, zero value and false otherwise.
	Get(ctx context.Context, key string) (T, bool)

	// Set stores a value in the cache with the specified TTL.
	// Returns an error if the operation fails.
	Set(ctx context.Context, key string, value T, ttl time.Duration) error

	// Delete removes a value from the cache.
	// Returns an error if the operation fails.
	Delete(ctx context.Context, key string) error

	// Close closes the cache connection/resources.
	// Should be called on application shutdown.
	Close() error
}

// HealthChecker is an optional interface that cache implementations
// can implement to provide health check functionality.
type HealthChecker interface {
	// Ping checks if the cache connection is alive.
	// Useful for health checks and readiness probes.
	Ping(ctx context.Context) error
}

// CacheType represents the type of cache implementation to use.
type CacheType string

const (
	// TypeMemory is an in-memory cache (non-distributed).
	TypeMemory CacheType = "memory"

	// TypeDistributed is a distributed cache backend.
	TypeDistributed CacheType = "distributed"

	// TypeNoOp is a no-op cache that does nothing (useful for testing).
	TypeNoOp CacheType = "noop"
)

// SerializationType represents the type of serialization to use.
type SerializationType string

const (
	// SerializationProtobuf uses Protocol Buffers for serialization.
	SerializationProtobuf SerializationType = "protobuf"
	// SerializationJSON uses JSON for serialization.
	SerializationJSON SerializationType = "json"
	// SerializationGob uses Go's gob encoding for serialization.
	SerializationGob SerializationType = "gob"
)
