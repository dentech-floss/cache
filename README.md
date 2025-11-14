# Cache

A flexible, type-safe caching library for Go with support for in-memory and distributed caching.

## Features

- **Type-safe** generic interface using Go generics
- **Multiple backends**: In-memory, distributed (Redis/Valkey), or no-op
- **Multiple serialization formats**: Protobuf, JSON, and Go binary (gob)
- **OpenTelemetry** instrumentation for observability
- **Health checks** for distributed backends
- **Simple API** with context support
- **Cloud Run optimized** for GCP services

## Under the Hood

This library builds on proven, battle-tested components:

- **In-memory cache**: [`github.com/jellydator/ttlcache/v2`](https://github.com/jellydator/ttlcache) - High-performance TTL cache
- **Distributed cache**: [`github.com/redis/go-redis/v9`](https://github.com/redis/go-redis) - Redis/Valkey client with connection pooling (supports both Redis and Valkey)
- **Protobuf support**: [`google.golang.org/protobuf`](https://pkg.go.dev/google.golang.org/protobuf) - Official protobuf library
- **Observability**: [`github.com/redis/go-redis/extra/redisotel/v9`](https://github.com/redis/go-redis/tree/master/extra/redisotel) - OpenTelemetry instrumentation

## Installation

```bash
go get github.com/dentech-floss/cache
```

## Quick Start

### Using the Factory (Recommended)

```go
import "github.com/dentech-floss/cache"

// Create a memory cache
config := &cache.Config{
    Type: cache.TypeMemory,
    Memory: &cache.MemoryConfig{
        SkipTTLExtensionOnHit: true,
    },
}

c, err := cache.New[*User](config)
if err != nil {
    panic(err)
}
defer c.Close()

// Use the cache
c.Set(ctx, "key", &User{ID: "123"}, 5*time.Minute)
user, found := c.Get(ctx, "key")
```

### Direct Creation

#### In-Memory Cache

```go
import "github.com/dentech-floss/cache"

// Create an in-memory cache for any type
c := cache.NewMemory[*User](nil)
defer c.Close()

// Use the cache
c.Set(ctx, "key", &User{ID: "123"}, 5*time.Minute)
user, found := c.Get(ctx, "key")
```

#### Distributed Cache (Protobuf)

```go
import "github.com/dentech-floss/cache"

// Create a distributed cache for protobuf messages
config := &cache.DistributedConfig{
    Addr: "localhost:6379",
}

c, err := cache.NewDistributed[*pb.User](config)
if err != nil {
    panic(err)
}
defer c.Close()

// Use the cache
c.Set(ctx, "key", &pb.User{Id: "123"}, 5*time.Minute)
user, found := c.Get(ctx, "key")
```

#### Distributed Cache (JSON/Generic)

```go
import "github.com/dentech-floss/cache"

// Create a distributed cache for any type with JSON serialization
config := &cache.DistributedConfig{
    Addr:              "localhost:6379",
    SerializationType: cache.SerializationJSON,
}

c, err := cache.NewDistributedGeneric[*User](config)
if err != nil {
    panic(err)
}
defer c.Close()

// Use the cache
c.Set(ctx, "key", &User{ID: "123"}, 5*time.Minute)
user, found := c.Get(ctx, "key")
```

#### Distributed Cache (Gob/Generic)

```go
import "github.com/dentech-floss/cache"

// Create a distributed cache for any type with Gob serialization (faster than JSON)
config := &cache.DistributedConfig{
    Addr:              "localhost:6379",
    SerializationType: cache.SerializationGob,
}

c, err := cache.NewDistributedGeneric[*User](config)
if err != nil {
    panic(err)
}
defer c.Close()

// Use the cache
c.Set(ctx, "key", &User{ID: "123"}, 5*time.Minute)
user, found := c.Get(ctx, "key")
```

#### No-Op Cache

```go
import "github.com/dentech-floss/cache"

// Useful for testing or when caching is disabled
c := cache.NewNoOp[*User]()
```

## Interface

All cache implementations satisfy the `Cache[T]` interface:

```go
type Cache[T any] interface {
    Get(ctx context.Context, key string) (T, bool)
    Set(ctx context.Context, key string, value T, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Close() error
}
```

## Configuration

### Memory Cache

```go
config := &cache.MemoryConfig{
    SkipTTLExtensionOnHit: true, // Don't extend TTL on cache hits
}
```

### Distributed Cache

```go
config := &cache.DistributedConfig{
    Addr:              "localhost:6379", // Works with both Redis and Valkey
    Password:          "optional-password",
    DB:                0,
    PoolSize:          10,
    MinIdleConns:      5,
    MaxRetries:        3,
    DialTimeout:       5 * time.Second,
    ReadTimeout:       3 * time.Second,
    WriteTimeout:      3 * time.Second,
    EnableTracing:     true,
    EnableMetrics:     true,
    SerializationType: cache.SerializationJSON, // or SerializationGob
    Client:            nil, // Optional: reuse an existing redis.UniversalClient
}
```

**Note**: The distributed cache works with both Redis and Valkey servers. Simply point the `Addr` to your Redis or Valkey instance.

### Sharing a Redis/Valkey Connection

You can now supply an existing `redis.UniversalClient` (for example, a `*redis.Client` or `*redis.ClusterClient`) so multiple caches reuse the same connection pool:

```go
import (
    "github.com/dentech-floss/cache"
    "github.com/redis/go-redis/v9"
)

sharedClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

userCache, err := cache.NewDistributedGeneric[*User](&cache.DistributedConfig{
    Client:            sharedClient,
    SerializationType: cache.SerializationJSON,
})
// handle error
orderCache, err := cache.NewDistributedGeneric[*Order](&cache.DistributedConfig{
    Client:            sharedClient,
    SerializationType: cache.SerializationGob,
})
// handle error
```

When a shared client is provided, the cache skips instrumentation and closing the client—allowing your application to manage its lifecycle centrally.

## Serialization Types

- **Protobuf**: For protobuf messages (automatic detection)
  - Best for: Microservices communication, when you already use protobuf
  - Pros: Compact, fast, schema evolution support
  - Cons: Requires protobuf definitions

- **JSON**: For any JSON-serializable type
  - Best for: General purpose, debugging, interoperability
  - Pros: Human-readable, language-agnostic, easy to debug
  - Cons: Larger size, slower than binary formats

- **Gob**: For any Go type (faster than JSON, but Go-specific)
  - Best for: Go-only environments, performance-critical applications
  - Pros: Fastest, smallest size, handles complex Go types
  - Cons: Go-specific, not human-readable

## Choosing the Right Cache Type

### Memory Cache (`TypeMemory`)
- **Use when**: Single instance, development, testing
- **Pros**: Fastest, no network overhead, simple setup
- **Cons**: Not shared between instances, lost on restart

### Distributed Cache (`TypeDistributed`)
- **Use when**: Multiple instances, production, shared state
- **Pros**: Shared between instances, persistent, scalable
- **Cons**: Network overhead, requires Redis/Valkey setup

### No-Op Cache (`TypeNoOp`)
- **Use when**: Testing, debugging, disabling cache
- **Pros**: No overhead, predictable behavior
- **Cons**: No caching benefits

## Health Checks

Distributed caches implement the `HealthChecker` interface:

```go
if healthChecker, ok := cache.(cache.HealthChecker); ok {
    err := healthChecker.Ping(ctx)
    if err != nil {
        // Handle unhealthy cache
    }
}
```

## Performance Considerations

- **Memory cache**: ~1-10μs per operation
- **Distributed cache**: ~100-1000μs per operation (network dependent)
- **Serialization overhead**: Gob < Protobuf < JSON
- **TTL precision**: Memory cache has second precision, distributed cache has millisecond precision

## Error Handling

The cache library follows Go's error handling conventions:

```go
// Set operations can fail
err := cache.Set(ctx, "key", value, ttl)
if err != nil {
    log.Printf("Cache set failed: %v", err)
    // Continue without caching
}

// Get operations return false on cache miss or error
value, found := cache.Get(ctx, "key")
if !found {
    // Cache miss - fetch from source
}

// Delete operations can fail
err := cache.Delete(ctx, "key")
if err != nil {
    log.Printf("Cache delete failed: %v", err)
}
```

## Best Practices

1. **Always handle errors** from Set/Delete operations
2. **Use context cancellation** for timeout control
3. **Choose appropriate TTL** based on your data freshness requirements
4. **Use NoOp cache** in tests for predictable behavior
5. **Monitor cache hit rates** and adjust TTL accordingly
6. **Use health checks** in production for distributed caches

## License

Apache 2.0 License
