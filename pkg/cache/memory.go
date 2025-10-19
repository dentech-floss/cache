package cache

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v2"
)

// memoryCache is an in-memory cache implementation.
type memoryCache[T any] struct {
	config *MemoryConfig
	cache  *ttlcache.Cache
}

// NewMemory creates a new in-memory cache with optional configuration.
// This is a convenience function for creating memory caches directly.
func NewMemory[T any](config *MemoryConfig) Cache[T] {
	cache := ttlcache.NewCache()

	if config != nil {
		cache.SkipTTLExtensionOnHit(config.SkipTTLExtensionOnHit)
	} else {
		// Default behavior: don't extend TTL on hit
		cache.SkipTTLExtensionOnHit(true)
	}

	return &memoryCache[T]{
		config: config,
		cache:  cache,
	}
}

func (c *memoryCache[T]) Get(ctx context.Context, key string) (T, bool) {
	var zero T

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return zero, false
	default:
	}

	if c.cache == nil {
		return zero, false
	}

	value, err := c.cache.Get(key)
	if err != nil {
		return zero, false
	}

	typedValue, ok := value.(T)
	if !ok {
		return zero, false
	}

	return typedValue, true
}

func (c *memoryCache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if c.cache == nil {
		return nil
	}

	return c.cache.SetWithTTL(key, value, ttl)
}

func (c *memoryCache[T]) Delete(ctx context.Context, key string) error {
	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if c.cache == nil {
		return nil
	}

	return c.cache.Remove(key)
}

func (c *memoryCache[T]) Close() error {
	if c.cache != nil {
		return c.cache.Close()
	}
	return nil
}
