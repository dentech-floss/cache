package cache

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

// distributedCache is a distributed cache implementation for proto messages.
type distributedCache[T any] struct {
	client *redis.Client
}

// distributedGenericCache is a distributed cache implementation for any type.
type distributedGenericCache[T any] struct {
	client     *redis.Client
	serializer Serializer
}

// NewDistributed creates a new distributed cache for proto messages.
// This is a convenience function for creating distributed caches directly.
func NewDistributed[T proto.Message](config *DistributedConfig) (Cache[T], error) {
	return NewDistributedForProto[T](config)
}

// NewDistributedForProto creates a new distributed cache for proto messages.
// This is an internal function used by the factory.
func NewDistributedForProto[T proto.Message](config *DistributedConfig) (Cache[T], error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	// Set defaults
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	if config.MinIdleConns == 0 {
		config.MinIdleConns = 5
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.DialTimeout == 0 {
		config.DialTimeout = 5 * time.Second
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 3 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 3 * time.Second
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	// Enable OpenTelemetry instrumentation
	if config.EnableTracing || config.EnableMetrics {
		if err := redisotel.InstrumentTracing(client); err != nil {
			return nil, err
		}
		if err := redisotel.InstrumentMetrics(client); err != nil {
			return nil, err
		}
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &distributedCache[T]{
		client: client,
	}, nil
}

// NewDistributedGeneric creates a new distributed cache for any type.
// This is a convenience function for creating distributed caches with custom serialization.
func NewDistributedGeneric[T any](config *DistributedConfig) (Cache[T], error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	// Set defaults
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	if config.MinIdleConns == 0 {
		config.MinIdleConns = 5
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.DialTimeout == 0 {
		config.DialTimeout = 5 * time.Second
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 3 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 3 * time.Second
	}

	// Set up serialization
	var serializer Serializer
	var err error

	if config.Serializer != nil {
		serializer = config.Serializer
	} else {
		// Default to JSON if no serialization type specified
		serializationType := config.SerializationType
		if serializationType == "" {
			serializationType = SerializationJSON
		}
		serializer, err = NewSerializer(serializationType)
		if err != nil {
			return nil, err
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	// Enable OpenTelemetry instrumentation
	if config.EnableTracing || config.EnableMetrics {
		if err := redisotel.InstrumentTracing(client); err != nil {
			return nil, err
		}
		if err := redisotel.InstrumentMetrics(client); err != nil {
			return nil, err
		}
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &distributedGenericCache[T]{
		client:     client,
		serializer: serializer,
	}, nil
}

// isProtoMessage checks if a type implements proto.Message using reflection
func isProtoMessage(v interface{}) bool {
	_, ok := v.(proto.Message)
	return ok
}

// createDistributedCacheForProto creates a distributed cache for proto messages
func createDistributedCacheForProto[T any](config *DistributedConfig) (Cache[T], error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	// Set defaults
	if config.PoolSize == 0 {
		config.PoolSize = 10
	}
	if config.MinIdleConns == 0 {
		config.MinIdleConns = 5
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.DialTimeout == 0 {
		config.DialTimeout = 5 * time.Second
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = 3 * time.Second
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = 3 * time.Second
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	// Enable OpenTelemetry instrumentation
	if config.EnableTracing || config.EnableMetrics {
		if err := redisotel.InstrumentTracing(client); err != nil {
			return nil, err
		}
		if err := redisotel.InstrumentMetrics(client); err != nil {
			return nil, err
		}
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), config.DialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &distributedCache[T]{
		client: client,
	}, nil
}

// Methods for distributedCache (proto messages)

func (c *distributedCache[T]) Get(ctx context.Context, key string) (T, bool) {
	var zero T

	if c.client == nil {
		return zero, false
	}

	// Get the serialized data
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		// Key not found or other error - treat as cache miss
		return zero, false
	}

	// Check if T is a proto.Message
	if _, ok := any(zero).(proto.Message); ok {
		// Create a new instance of T using reflection
		result := reflect.New(reflect.TypeOf(zero).Elem()).Interface().(T)

		// Deserialize the proto message
		if err := proto.Unmarshal(data, any(result).(proto.Message)); err != nil {
			// Failed to deserialize - treat as cache miss
			return zero, false
		}

		return result, true
	}

	// This should not happen if we're using this cache correctly
	return zero, false
}

func (c *distributedCache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}

	// Check if T is a proto.Message
	if protoMsg, ok := any(value).(proto.Message); ok {
		// Serialize the proto message
		data, err := proto.Marshal(protoMsg)
		if err != nil {
			return err
		}

		// Store with TTL
		return c.client.Set(ctx, key, data, ttl).Err()
	}

	// This should not happen if we're using this cache correctly
	return errors.New("distributedCache can only be used with proto.Message types")
}

func (c *distributedCache[T]) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return nil
	}

	return c.client.Del(ctx, key).Err()
}

func (c *distributedCache[T]) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *distributedCache[T]) Ping(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Ping(ctx).Err()
}

// Methods for distributedGenericCache (any type)

func (c *distributedGenericCache[T]) Get(ctx context.Context, key string) (T, bool) {
	var zero T

	if c.client == nil {
		return zero, false
	}

	// Get the serialized data
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		// Key not found or other error - treat as cache miss
		return zero, false
	}

	// Create a new instance of T
	var result T

	// Deserialize the data
	if err := c.serializer.Deserialize(data, &result); err != nil {
		// Failed to deserialize - treat as cache miss
		return zero, false
	}

	return result, true
}

func (c *distributedGenericCache[T]) Set(ctx context.Context, key string, value T, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}

	// Serialize the value
	data, err := c.serializer.Serialize(value)
	if err != nil {
		return err
	}

	// Store with TTL
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *distributedGenericCache[T]) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return nil
	}

	return c.client.Del(ctx, key).Err()
}

func (c *distributedGenericCache[T]) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

func (c *distributedGenericCache[T]) Ping(ctx context.Context) error {
	if c.client == nil {
		return nil
	}
	return c.client.Ping(ctx).Err()
}
