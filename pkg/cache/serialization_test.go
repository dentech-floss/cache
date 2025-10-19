package cache

import (
	"testing"
)

func TestSerializers(t *testing.T) {
	// Test JSON Serializer
	jsonSerializer := &JSONSerializer{}
	user := TestUser{ID: "123", Name: "John"}

	data, err := jsonSerializer.Serialize(user)
	if err != nil {
		t.Errorf("JSON Serialize failed: %v", err)
	}

	var retrieved TestUser
	err = jsonSerializer.Deserialize(data, &retrieved)
	if err != nil {
		t.Errorf("JSON Deserialize failed: %v", err)
	}

	if retrieved.ID != user.ID || retrieved.Name != user.Name {
		t.Errorf("Expected %+v, got %+v", user, retrieved)
	}

	// Test Gob Serializer
	gobSerializer := &GobSerializer{}

	data, err = gobSerializer.Serialize(user)
	if err != nil {
		t.Errorf("Gob Serialize failed: %v", err)
	}

	var retrieved2 TestUser
	err = gobSerializer.Deserialize(data, &retrieved2)
	if err != nil {
		t.Errorf("Gob Deserialize failed: %v", err)
	}

	if retrieved2.ID != user.ID || retrieved2.Name != user.Name {
		t.Errorf("Expected %+v, got %+v", user, retrieved2)
	}
}

func TestNewSerializer(t *testing.T) {
	tests := []struct {
		name              string
		serializationType SerializationType
		wantErr           bool
		errorMsg          string
	}{
		{
			name:              "JSON serialization",
			serializationType: SerializationJSON,
			wantErr:           false,
		},
		{
			name:              "Gob serialization",
			serializationType: SerializationGob,
			wantErr:           false,
		},
		{
			name:              "Protobuf serialization",
			serializationType: SerializationProtobuf,
			wantErr:           true,
			errorMsg:          "protobuf serialization requires special handling - use NewDistributed",
		},
		{
			name:              "Unknown serialization",
			serializationType: SerializationType("unknown"),
			wantErr:           true,
			errorMsg:          "unknown serialization type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serializer, err := NewSerializer(tt.serializationType)
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
			if serializer == nil {
				t.Error("Expected serializer, got nil")
			}
		})
	}
}
