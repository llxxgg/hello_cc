package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateUUID(t *testing.T) {
	// Test that GenerateUUID returns a non-empty string
	uuidStr := GenerateUUID()
	if uuidStr == "" {
		t.Error("GenerateUUID() should not return empty string")
	}

	// Test that the returned string is a valid UUID format (with hyphens)
	_, err := uuid.Parse(uuidStr)
	if err != nil {
		t.Errorf("GenerateUUID() should return valid UUID, got: %s, error: %v", uuidStr, err)
	}

	// Test that GenerateUUID returns different UUIDs each time
	uuid1 := GenerateUUID()
	uuid2 := GenerateUUID()
	if uuid1 == uuid2 {
		t.Error("GenerateUUID() should return different UUIDs on each call")
	}

	// Test UUID format contains expected hyphens (8-4-4-4-12 format)
	expectedLen := 36 // xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if len(uuidStr) != expectedLen {
		t.Errorf("GenerateUUID() should return UUID of length %d, got: %d", expectedLen, len(uuidStr))
	}
}

func TestGenerateUUIDShort(t *testing.T) {
	// Test that GenerateUUIDShort returns a non-empty string
	uuidStr := GenerateUUIDShort()
	if uuidStr == "" {
		t.Error("GenerateUUIDShort() should not return empty string")
	}

	// Test that the returned string is a valid UUID format
	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		t.Errorf("GenerateUUIDShort() should return valid UUID, got: %s, error: %v", uuidStr, err)
	}

	// Test that GenerateUUIDShort returns different UUIDs each time
	uuid1 := GenerateUUIDShort()
	uuid2 := GenerateUUIDShort()
	if uuid1 == uuid2 {
		t.Error("GenerateUUIDShort() should return different UUIDs on each call")
	}

	// Test that the UUID is variant 2 (DCE 1.1, ISO/IEC 11578)
	if parsedUUID.Variant() != uuid.RFC4122 {
		t.Errorf("GenerateUUIDShort() should return RFC4122 variant UUID, got: %v", parsedUUID.Variant())
	}

	// Test that the UUID is version 3 (MD5)
	if parsedUUID.Version() != 3 {
		t.Errorf("GenerateUUIDShort() should return MD5 version UUID, got: %v", parsedUUID.Version())
	}

	// Test UUID format (without hyphens, so 32 characters)
	expectedLen := 32
	if len(uuidStr) != expectedLen {
		t.Errorf("GenerateUUIDShort() should return UUID of length %d, got: %d", expectedLen, len(uuidStr))
	}
}
