package cache

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds common configuration for all cache types.
type Config struct {
	// Type specifies which cache implementation to use.
	Type CacheType

	// Memory-specific configuration (only used when Type is TypeMemory)
	Memory *MemoryConfig

	// Distributed-specific configuration (only used when Type is TypeDistributed)
	Distributed *DistributedConfig
}

// MemoryConfig holds configuration for in-memory cache.
type MemoryConfig struct {
	// SkipTTLExtensionOnHit prevents TTL from being reset on cache hits.
	// Default: true
	SkipTTLExtensionOnHit bool
}

// DistributedConfig holds configuration for distributed cache.
type DistributedConfig struct {
	// Addr is the cache server address (e.g., "localhost:6379")
	Addr string

	// Password for authentication (optional)
	Password string

	// DB is the database number to use (default: 0)
	DB int

	// PoolSize is the maximum number of socket connections (default: 10)
	PoolSize int

	// MinIdleConns is the minimum number of idle connections (default: 5)
	MinIdleConns int

	// MaxRetries is the maximum number of retries before giving up (default: 3)
	MaxRetries int

	// DialTimeout is the timeout for establishing new connections (default: 5s)
	DialTimeout time.Duration

	// ReadTimeout is the timeout for socket reads (default: 3s)
	ReadTimeout time.Duration

	// WriteTimeout is the timeout for socket writes (default: 3s)
	WriteTimeout time.Duration

	// EnableTracing enables OpenTelemetry tracing for cache operations (default: true)
	EnableTracing bool

	// EnableMetrics enables OpenTelemetry metrics for cache operations (default: true)
	EnableMetrics bool

	// SerializationType specifies how to serialize data (default: protobuf for proto.Message, json for others)
	SerializationType SerializationType

	// Serializer allows custom serialization (overrides SerializationType if set)
	Serializer Serializer

	// Client allows providing a pre-configured Redis/Valkey client.
	// When set, the cache will reuse this client instead of creating its own.
	// The cache will not close the shared client when Close is called.
	Client redis.UniversalClient
}
