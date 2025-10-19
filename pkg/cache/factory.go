package cache

import (
	"errors"
	"fmt"
)

// New creates a new cache based on the provided configuration.
// This is the recommended way to create caches as it handles all the setup.
func New[T any](config *Config) (Cache[T], error) {
	if config == nil {
		return nil, errors.New("config cannot be nil")
	}

	switch config.Type {
	case TypeMemory:
		return NewMemory[T](config.Memory), nil

	case TypeDistributed:
		// For distributed cache, we need to check if T is a proto.Message
		var zero T
		if isProtoMessage(zero) {
			// Use the protobuf-specific implementation
			return createDistributedCacheForProto[T](config.Distributed)
		} else {
			// Use the generic implementation for non-proto types
			return NewDistributedGeneric[T](config.Distributed)
		}

	case TypeNoOp:
		return NewNoOp[T](), nil

	default:
		return nil, fmt.Errorf("unknown cache type: %s", config.Type)
	}
}
