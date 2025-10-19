package cache

import (
	"testing"
)

func TestFactoryNew(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		wantErr  bool
		errorMsg string
	}{
		{
			name:     "nil config",
			config:   nil,
			wantErr:  true,
			errorMsg: "config cannot be nil",
		},
		{
			name: "memory cache",
			config: &Config{
				Type: TypeMemory,
			},
			wantErr: false,
		},
		{
			name: "no-op cache",
			config: &Config{
				Type: TypeNoOp,
			},
			wantErr: false,
		},
		{
			name: "unknown cache type",
			config: &Config{
				Type: CacheType("unknown"),
			},
			wantErr:  true,
			errorMsg: "unknown cache type: unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache, err := New[TestUser](tt.config)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if cache == nil {
				t.Error("Expected cache, got nil")
			}
		})
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Test NewMemory
	cache1 := NewMemory[TestUser](nil)
	if cache1 == nil {
		t.Error("NewMemory returned nil")
	}
	_ = cache1.Close()

	// Test NewNoOp
	cache2 := NewNoOp[TestUser]()
	if cache2 == nil {
		t.Error("NewNoOp returned nil")
	}
	_ = cache2.Close()

	// Test NewMemory with config
	config := &MemoryConfig{SkipTTLExtensionOnHit: true}
	cache3 := NewMemory[TestUser](config)
	if cache3 == nil {
		t.Error("NewMemory with config returned nil")
	}
	_ = cache3.Close()
}
