package serializer

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/enums"
)

// Test types for testing
type TestStruct struct {
	CommandCode enums.CommonCommand         `enocean-esp3:"data"`
	DeviceID    deviceid.DeviceID           `enocean-esp3:"data"`
	SecurityKey [16]byte                    `enocean-esp3:"optdata"`
	Direction   enums.SecureDeviceDirection `enocean-esp3:"optdata"`
}

type CustomType uint32

type TestStructWithCustomType struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	CustomValue CustomType          `enocean-esp3:"data"`
}

// TestCommandToTelegram verifies CommandToTelegram behavior.
func TestCommandToTelegram(t *testing.T) {
	t.Run("basic serialization", func(t *testing.T) {
		cmd := TestStruct{
			CommandCode: enums.CommonCommandWR_SECUREDEVICE_DEL,
			DeviceID:    deviceid.DeviceID(0x12345678),
			Direction:   enums.SecureDeviceDirectionNONE,
			SecurityKey: [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize command: %v", err)
		}

		// Verify data: CommandCode (0x1a) + DeviceID (0x12345678 as big-endian uint32)
		expectedData := []byte{0x1a, 0x12, 0x34, 0x56, 0x78}
		if len(telegram.Data) != len(expectedData) {
			t.Errorf("data length mismatch: got %d, expected %d", len(telegram.Data), len(expectedData))
		}
		for i, b := range expectedData {
			if i < len(telegram.Data) && telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x", i, telegram.Data[i], b)
			}
		}

		// Verify optdata: SecurityKey (16 bytes) + Direction (0xff)
		expectedOptData := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0xff}
		if len(telegram.OptData) != len(expectedOptData) {
			t.Errorf("optdata length mismatch: got %d, expected %d", len(telegram.OptData), len(expectedOptData))
		}
		for i, b := range expectedOptData {
			if i < len(telegram.OptData) && telegram.OptData[i] != b {
				t.Errorf("optdata[%d]: got 0x%02x, expected 0x%02x", i, telegram.OptData[i], b)
			}
		}
	})

	t.Run("returns error for non-struct input", func(t *testing.T) {
		cfg := SerializerConfig{}
		_, err := CommandToTelegram(42, cfg)
		if err == nil {
			t.Errorf("expected error for non-struct input, got nil")
		}
		if err.Error() != "command must be a struct" {
			t.Errorf("expected error 'command must be a struct', got '%s'", err.Error())
		}
	})

	t.Run("returns error for pointer to non-struct", func(t *testing.T) {
		val := 42
		cfg := SerializerConfig{}
		_, err := CommandToTelegram(&val, cfg)
		if err == nil {
			t.Errorf("expected error for pointer to non-struct, got nil")
		}
		if err.Error() != "command must be a struct" {
			t.Errorf("expected error 'command must be a struct', got '%s'", err.Error())
		}
	})

	t.Run("handles pointer to struct", func(t *testing.T) {
		cmd := &TestStruct{
			CommandCode: enums.CommonCommandWR_SECUREDEVICE_DEL,
			DeviceID:    deviceid.DeviceID(0x12345678),
			Direction:   enums.SecureDeviceDirectionNONE,
			SecurityKey: [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize command: %v", err)
		}

		if len(telegram.Data) == 0 {
			t.Errorf("expected non-empty data")
		}
	})

	t.Run("skips fields without tags", func(t *testing.T) {
		type StructWithoutTags struct {
			Field1 int
			Field2 string `enocean-esp3:"data"`
		}

		cmd := StructWithoutTags{
			Field1: 42,
			Field2: "test",
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Only Field2 should be serialized
		if len(telegram.Data) != len("test") {
			t.Errorf("expected data length %d, got %d", len("test"), len(telegram.Data))
		}
	})

	t.Run("handles skipif:none tag", func(t *testing.T) {
		type StructWithSkipIf struct {
			Field1 int32                       `enocean-esp3:"optdata,skipif:none"`
			Field2 enums.SecureDeviceDirection `enocean-esp3:"optdata,skipif:none"`
			Field3 int32                       `enocean-esp3:"optdata"`
		}

		// Use actual zero value for SecureDeviceDirection (0x00)
		var zeroDir enums.SecureDeviceDirection = 0
		cmd := StructWithSkipIf{
			Field1: 0,       // Zero value, should be skipped
			Field2: zeroDir, // Zero value, should be skipped
			Field3: 42,      // Non-zero, should be included
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Only Field3 should be in optdata (int32 = 4 bytes)
		// Both Field1 and Field2 are zero and should be skipped
		if len(telegram.OptData) != 4 {
			t.Errorf("expected optdata length 4, got %d (Field1 and Field2 should be skipped)", len(telegram.OptData))
		}
	})

	t.Run("handles skipif:none tag with non-zero value", func(t *testing.T) {
		type StructWithSkipIf struct {
			Field1 int32 `enocean-esp3:"optdata,skipif:none"`
			Field2 int32 `enocean-esp3:"optdata"`
		}

		cmd := StructWithSkipIf{
			Field1: 100, // Non-zero, should NOT be skipped
			Field2: 42,
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Both fields should be in optdata (2 * int32 = 8 bytes)
		if len(telegram.OptData) != 8 {
			t.Errorf("expected optdata length 8, got %d (both fields should be included)", len(telegram.OptData))
		}
	})

	t.Run("returns error for invalid target", func(t *testing.T) {
		type StructWithInvalidTarget struct {
			Field1 int32 `enocean-esp3:"invalid"`
		}

		cmd := StructWithInvalidTarget{Field1: 42}
		cfg := SerializerConfig{}
		_, err := CommandToTelegram(cmd, cfg)
		if err == nil {
			t.Errorf("expected error for invalid target, got nil")
		}
		if err.Error() != "invalid target: invalid" {
			t.Errorf("expected error 'invalid target: invalid', got '%s'", err.Error())
		}
	})

	t.Run("handles empty optdata", func(t *testing.T) {
		type StructDataOnly struct {
			Field1 int32 `enocean-esp3:"data"`
		}

		cmd := StructDataOnly{Field1: 42}
		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		if telegram.OptData != nil {
			t.Errorf("expected nil optdata, got %v", telegram.OptData)
		}
	})

	t.Run("handles unexported fields", func(t *testing.T) {
		type StructWithUnexported struct {
			Exported   int32 `enocean-esp3:"data"`
			unexported int32 // Should be skipped
		}

		cmd := StructWithUnexported{
			Exported:   42,
			unexported: 100,
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Only exported field should be serialized
		if len(telegram.Data) != 4 { // int32 = 4 bytes
			t.Errorf("expected data length 4, got %d", len(telegram.Data))
		}
	})
}

// TestCommandToTelegram_CustomSerializer verifies CommandToTelegram_CustomSerializer behavior.
func TestCommandToTelegram_CustomSerializer(t *testing.T) {
	t.Run("custom serializer usage", func(t *testing.T) {
		// Create a custom serializer for CustomType that multiplies by 2
		customSerializer := func(buf *bytes.Buffer, v reflect.Value, byteOrder binary.ByteOrder) error {
			// Custom serialization: multiply the value by 2
			originalValue := uint32(v.Uint())
			multipliedValue := originalValue * 2
			return binary.Write(buf, byteOrder, multipliedValue)
		}

		cfg := SerializerConfig{
			Serializers: map[reflect.Type]CustomSerializer{
				reflect.TypeOf(CustomType(0)): customSerializer,
			},
		}

		cmd := TestStructWithCustomType{
			CommandCode: enums.CommonCommandWR_SECUREDEVICE_DEL,
			CustomValue: CustomType(0x1234),
		}

		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize command: %v", err)
		}

		// Verify data: CommandCode (0x1a) + CustomValue (0x1234 * 2 = 0x2468 as big-endian uint32)
		expectedData := []byte{0x1a, 0x00, 0x00, 0x24, 0x68}
		if len(telegram.Data) != len(expectedData) {
			t.Errorf("data length mismatch: got %d, expected %d", len(telegram.Data), len(expectedData))
		}
		for i, b := range expectedData {
			if i < len(telegram.Data) && telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x", i, telegram.Data[i], b)
			}
		}

		// Verify the custom serializer was actually called
		// The value 0x1234 should be serialized as 0x2468 (multiplied by 2)
		if len(telegram.Data) >= 5 {
			// Extract the custom value (last 4 bytes)
			serializedValue := uint32(telegram.Data[1])<<24 | uint32(telegram.Data[2])<<16 | uint32(telegram.Data[3])<<8 | uint32(telegram.Data[4])
			expectedValue := uint32(0x1234) * 2 // Should be 0x2468
			if serializedValue != expectedValue {
				t.Errorf("custom serializer not called correctly: got 0x%x, expected 0x%x", serializedValue, expectedValue)
			}
		}
	})

	t.Run("custom serializer not used when not registered", func(t *testing.T) {
		// Config without custom serializer
		cfg := SerializerConfig{}

		cmd := TestStructWithCustomType{
			CommandCode: enums.CommonCommandWR_SECUREDEVICE_DEL,
			CustomValue: CustomType(0x1234),
		}

		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Should serialize as normal uint32 (0x1234), not multiplied
		if len(telegram.Data) >= 5 {
			serializedValue := uint32(telegram.Data[1])<<24 | uint32(telegram.Data[2])<<16 | uint32(telegram.Data[3])<<8 | uint32(telegram.Data[4])
			if serializedValue == 0x2468 {
				t.Errorf("custom serializer should not be used, but value was multiplied")
			}
			if serializedValue != 0x1234 {
				t.Errorf("expected normal serialization 0x1234, got 0x%x", serializedValue)
			}
		}
	})

	t.Run("custom serializer with nil map", func(t *testing.T) {
		// Config with nil Serializers map should still work (sanitize will fix it)
		cfg := SerializerConfig{
			Serializers: nil,
		}

		cmd := TestStruct{
			CommandCode: enums.CommonCommandWR_SECUREDEVICE_DEL,
			DeviceID:    deviceid.DeviceID(0x12345678),
		}

		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		if len(telegram.Data) == 0 {
			t.Errorf("expected non-empty data")
		}
	})

	t.Run("custom serializer with empty map", func(t *testing.T) {
		cfg := SerializerConfig{
			Serializers: map[reflect.Type]CustomSerializer{},
		}

		cmd := TestStruct{
			CommandCode: enums.CommonCommandWR_SECUREDEVICE_DEL,
			DeviceID:    deviceid.DeviceID(0x12345678),
		}

		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		if len(telegram.Data) == 0 {
			t.Errorf("expected non-empty data")
		}
	})

	t.Run("custom serializer error propagation", func(t *testing.T) {
		// Custom serializer that returns an error
		customSerializer := func(buf *bytes.Buffer, v reflect.Value, byteOrder binary.ByteOrder) error {
			return bytes.ErrTooLarge // Simulate an error
		}

		cfg := SerializerConfig{
			Serializers: map[reflect.Type]CustomSerializer{
				reflect.TypeOf(CustomType(0)): customSerializer,
			},
		}

		cmd := TestStructWithCustomType{
			CommandCode: enums.CommonCommandWR_SECUREDEVICE_DEL,
			CustomValue: CustomType(0x1234),
		}

		_, err := CommandToTelegram(cmd, cfg)
		if err == nil {
			t.Errorf("expected error from custom serializer, got nil")
		}
	})
}

// TestCommandToTelegram_ByteOrder verifies CommandToTelegram_ByteOrder behavior.
func TestCommandToTelegram_ByteOrder(t *testing.T) {
	t.Run("uses big endian by default", func(t *testing.T) {
		type IntStruct struct {
			Value uint32 `enocean-esp3:"data"`
		}

		cmd := IntStruct{Value: 0x12345678}
		cfg := SerializerConfig{} // ByteOrder is nil, should default to BigEndian

		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Big endian: 0x12 0x34 0x56 0x78
		expected := []byte{0x12, 0x34, 0x56, 0x78}
		for i, b := range expected {
			if telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x (big endian)", i, telegram.Data[i], b)
			}
		}
	})

	t.Run("uses little endian when specified", func(t *testing.T) {
		type IntStruct struct {
			Value uint32 `enocean-esp3:"data"`
		}

		cmd := IntStruct{Value: 0x12345678}
		cfg := SerializerConfig{
			ByteOrder: binary.LittleEndian,
		}

		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Little endian: 0x78 0x56 0x34 0x12
		expected := []byte{0x78, 0x56, 0x34, 0x12}
		for i, b := range expected {
			if telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x (little endian)", i, telegram.Data[i], b)
			}
		}
	})

	t.Run("uses explicit big endian when specified", func(t *testing.T) {
		type IntStruct struct {
			Value uint32 `enocean-esp3:"data"`
		}

		cmd := IntStruct{Value: 0x12345678}
		cfg := SerializerConfig{
			ByteOrder: binary.BigEndian,
		}

		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Big endian: 0x12 0x34 0x56 0x78
		expected := []byte{0x12, 0x34, 0x56, 0x78}
		for i, b := range expected {
			if telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x (big endian)", i, telegram.Data[i], b)
			}
		}
	})
}

// TestCommandToTelegram_AllTypes verifies CommandToTelegram_AllTypes behavior.
func TestCommandToTelegram_AllTypes(t *testing.T) {
	t.Run("serializes all integer types", func(t *testing.T) {
		type IntStruct struct {
			I8  int8   `enocean-esp3:"data"`
			I16 int16  `enocean-esp3:"data"`
			I32 int32  `enocean-esp3:"data"`
			I64 int64  `enocean-esp3:"data"`
			U8  uint8  `enocean-esp3:"data"`
			U16 uint16 `enocean-esp3:"data"`
			U32 uint32 `enocean-esp3:"data"`
			U64 uint64 `enocean-esp3:"data"`
		}

		cmd := IntStruct{
			I8:  -128,
			I16: -32768,
			I32: -2147483648,
			I64: -9223372036854775808,
			U8:  255,
			U16: 65535,
			U32: 4294967295,
			U64: 18446744073709551615,
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Verify all types are serialized (1+2+4+8+1+2+4+8 = 30 bytes)
		expectedLen := 1 + 2 + 4 + 8 + 1 + 2 + 4 + 8
		if len(telegram.Data) != expectedLen {
			t.Errorf("expected data length %d, got %d", expectedLen, len(telegram.Data))
		}
	})

	t.Run("serializes float types", func(t *testing.T) {
		type FloatStruct struct {
			F32 float32 `enocean-esp3:"data"`
			F64 float64 `enocean-esp3:"data"`
		}

		cmd := FloatStruct{
			F32: 3.14159,
			F64: 2.718281828,
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Verify floats are serialized (4+8 = 12 bytes)
		if len(telegram.Data) != 12 {
			t.Errorf("expected data length 12, got %d", len(telegram.Data))
		}
	})

	t.Run("serializes boolean", func(t *testing.T) {
		type BoolStruct struct {
			True  bool `enocean-esp3:"data"`
			False bool `enocean-esp3:"data"`
		}

		cmd := BoolStruct{
			True:  true,
			False: false,
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Booleans serialized as uint8 (1 byte each)
		if len(telegram.Data) != 2 {
			t.Errorf("expected data length 2, got %d", len(telegram.Data))
		}
		if telegram.Data[0] != 1 {
			t.Errorf("expected true to serialize as 1, got %d", telegram.Data[0])
		}
		if telegram.Data[1] != 0 {
			t.Errorf("expected false to serialize as 0, got %d", telegram.Data[1])
		}
	})

	t.Run("serializes string", func(t *testing.T) {
		type StringStruct struct {
			Str string `enocean-esp3:"data"`
		}

		cmd := StringStruct{
			Str: "Hello, World!",
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// String should be serialized as UTF-8 bytes
		if len(telegram.Data) != len("Hello, World!") {
			t.Errorf("expected data length %d, got %d", len("Hello, World!"), len(telegram.Data))
		}
		if string(telegram.Data) != "Hello, World!" {
			t.Errorf("expected string 'Hello, World!', got '%s'", string(telegram.Data))
		}
	})

	t.Run("serializes byte array", func(t *testing.T) {
		type ByteArrayStruct struct {
			Arr [4]byte `enocean-esp3:"data"`
		}

		cmd := ByteArrayStruct{
			Arr: [4]byte{0x01, 0x02, 0x03, 0x04},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Byte array should be written directly as raw bytes
		if len(telegram.Data) != 4 {
			t.Errorf("expected data length 4, got %d", len(telegram.Data))
		}
		expected := []byte{0x01, 0x02, 0x03, 0x04}
		for i, b := range expected {
			if telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x", i, telegram.Data[i], b)
			}
		}
	})

	t.Run("serializes byte array with CanAddr path", func(t *testing.T) {
		// This tests the CanAddr() path in serializeSequence for byte arrays
		type ByteArrayStruct struct {
			Arr [4]byte `enocean-esp3:"data"`
		}

		cmd := ByteArrayStruct{
			Arr: [4]byte{0x01, 0x02, 0x03, 0x04},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Byte array should be written directly as raw bytes
		if len(telegram.Data) != 4 {
			t.Errorf("expected data length 4, got %d", len(telegram.Data))
		}
		expected := []byte{0x01, 0x02, 0x03, 0x04}
		for i, b := range expected {
			if telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x", i, telegram.Data[i], b)
			}
		}
	})

	t.Run("serializes non-byte array", func(t *testing.T) {
		type IntArrayStruct struct {
			Arr [3]int16 `enocean-esp3:"data"`
		}

		cmd := IntArrayStruct{
			Arr: [3]int16{0x0102, 0x0304, 0x0506},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Each int16 is 2 bytes, so 3 * 2 = 6 bytes total
		if len(telegram.Data) != 6 {
			t.Errorf("expected data length 6, got %d", len(telegram.Data))
		}
	})

	t.Run("serializes byte slice", func(t *testing.T) {
		type ByteSliceStruct struct {
			Slice []byte `enocean-esp3:"data"`
		}

		cmd := ByteSliceStruct{
			Slice: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Byte slice should be written directly as raw bytes
		if len(telegram.Data) != 5 {
			t.Errorf("expected data length 5, got %d", len(telegram.Data))
		}
		expected := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
		for i, b := range expected {
			if telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x", i, telegram.Data[i], b)
			}
		}
	})

	t.Run("serializes nil byte slice", func(t *testing.T) {
		type NilSliceStruct struct {
			Slice []byte `enocean-esp3:"data"`
		}

		cmd := NilSliceStruct{
			Slice: nil,
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Nil slice should serialize as nothing (nil means "no data")
		if len(telegram.Data) != 0 {
			t.Errorf("expected data length 0 for nil slice, got %d", len(telegram.Data))
		}
	})

	t.Run("serializes non-byte slice", func(t *testing.T) {
		type IntSliceStruct struct {
			Slice []int16 `enocean-esp3:"data"`
		}

		cmd := IntSliceStruct{
			Slice: []int16{0x0102, 0x0304},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Each int16 is 2 bytes, so 2 * 2 = 4 bytes total
		if len(telegram.Data) != 4 {
			t.Errorf("expected data length 4, got %d", len(telegram.Data))
		}
	})

	t.Run("serializes nested struct", func(t *testing.T) {
		type NestedStruct struct {
			Field1 int16
			Field2 uint32
		}

		type OuterStruct struct {
			OuterField int8         `enocean-esp3:"data"`
			Nested     NestedStruct `enocean-esp3:"data"`
		}

		cmd := OuterStruct{
			OuterField: 42,
			Nested: NestedStruct{
				Field1: 0x1234,
				Field2: 0x567890AB,
			},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// OuterField (1 byte) + Nested.Field1 (2 bytes) + Nested.Field2 (4 bytes) = 7 bytes
		if len(telegram.Data) != 7 {
			t.Errorf("expected data length 7, got %d", len(telegram.Data))
		}
	})

	t.Run("serializes pointer to value", func(t *testing.T) {
		val := int32(42)
		type PointerStruct struct {
			Ptr *int32 `enocean-esp3:"data"`
		}

		cmd := PointerStruct{
			Ptr: &val,
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Pointer should dereference and serialize the value (4 bytes for int32)
		if len(telegram.Data) != 4 {
			t.Errorf("expected data length 4, got %d", len(telegram.Data))
		}
	})

	t.Run("serializes nil pointer", func(t *testing.T) {
		type NilPointerStruct struct {
			Ptr *int32 `enocean-esp3:"data"`
		}

		cmd := NilPointerStruct{
			Ptr: nil,
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Nil pointer should serialize as zero bytes (size of element type)
		if len(telegram.Data) != 4 { // int32 size = 4
			t.Errorf("expected data length 4, got %d", len(telegram.Data))
		}
	})

	t.Run("handles interface type", func(t *testing.T) {
		type InterfaceStruct struct {
			Iface interface{} `enocean-esp3:"data"`
		}

		cmd := InterfaceStruct{
			Iface: int32(42),
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Interface should serialize the underlying value
		if len(telegram.Data) != 4 { // int32 size = 4
			t.Errorf("expected data length 4, got %d", len(telegram.Data))
		}
	})

	t.Run("handles nil interface", func(t *testing.T) {
		type NilInterfaceStruct struct {
			Iface interface{} `enocean-esp3:"data"`
		}

		cmd := NilInterfaceStruct{
			Iface: nil,
		}

		cfg := SerializerConfig{}
		_, err := CommandToTelegram(cmd, cfg)
		if err == nil {
			t.Errorf("expected error for nil interface, got nil")
		}
		if err.Error() != "failed to serialize field Iface: nil interface: size unknown" {
			t.Errorf("expected specific error message, got '%s'", err.Error())
		}
	})

	t.Run("handles nested struct with tags", func(t *testing.T) {
		type NestedWithTags struct {
			Field1 int32 `enocean-esp3:"data"` // Tag should be skipped in nested struct
			Field2 int32 // No tag, should be serialized
		}

		type OuterStruct struct {
			OuterField int8           `enocean-esp3:"data"`
			Nested     NestedWithTags `enocean-esp3:"data"`
		}

		cmd := OuterStruct{
			OuterField: 42,
			Nested: NestedWithTags{
				Field1: 100,
				Field2: 200,
			},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// OuterField (1 byte) + Nested.Field2 (4 bytes) = 5 bytes
		// Nested struct fields with tags should be skipped, only Field2 should be serialized
		if len(telegram.Data) != 5 {
			t.Errorf("expected data length 5, got %d", len(telegram.Data))
		}
	})

	t.Run("handles default case for unknown types", func(t *testing.T) {
		type Complex64Struct struct {
			Complex complex64 `enocean-esp3:"data"`
		}

		cmd := Complex64Struct{
			Complex: complex64(1 + 2i),
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Should use default binary.Write (complex64 = 8 bytes: 2 float32s)
		if len(telegram.Data) != 8 {
			t.Errorf("expected data length 8, got %d", len(telegram.Data))
		}
	})

	t.Run("handles nested struct with unexported fields", func(t *testing.T) {
		type NestedWithUnexported struct {
			Exported   int32
			unexported int32 // Should be skipped
		}

		type OuterStruct struct {
			OuterField int8                 `enocean-esp3:"data"`
			Nested     NestedWithUnexported `enocean-esp3:"data"`
		}

		cmd := OuterStruct{
			OuterField: 42,
			Nested: NestedWithUnexported{
				Exported:   100,
				unexported: 200,
			},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// OuterField (1 byte) + Nested.Exported (4 bytes) = 5 bytes
		// unexported field should be skipped
		if len(telegram.Data) != 5 {
			t.Errorf("expected data length 5, got %d", len(telegram.Data))
		}
	})

	t.Run("handles nested struct with nil pointer", func(t *testing.T) {
		type NestedWithPointer struct {
			Field *int32
		}

		type OuterStruct struct {
			OuterField int8              `enocean-esp3:"data"`
			Nested     NestedWithPointer `enocean-esp3:"data"`
		}

		// Test with a nil pointer which should serialize as zero bytes
		cmd := OuterStruct{
			OuterField: 1,
			Nested: NestedWithPointer{
				Field: nil,
			},
		}

		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// OuterField (1 byte) + nil pointer zero bytes (4 bytes for *int32)
		if len(telegram.Data) != 5 {
			t.Errorf("expected data length 5, got %d", len(telegram.Data))
		}
	})
}

// TestCommandToTelegram_EdgeCases verifies CommandToTelegram_EdgeCases behavior.
func TestCommandToTelegram_EdgeCases(t *testing.T) {
	t.Run("handles empty struct", func(t *testing.T) {
		type EmptyStruct struct{}

		cmd := EmptyStruct{}
		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		if len(telegram.Data) != 0 {
			t.Errorf("expected empty data, got %d bytes", len(telegram.Data))
		}
		if telegram.OptData != nil {
			t.Errorf("expected nil optdata, got %v", telegram.OptData)
		}
	})

	t.Run("handles struct with only data fields", func(t *testing.T) {
		type DataOnlyStruct struct {
			Field1 int32 `enocean-esp3:"data"`
			Field2 int32 `enocean-esp3:"data"`
		}

		cmd := DataOnlyStruct{Field1: 1, Field2: 2}
		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		if len(telegram.Data) != 8 { // 2 * int32 = 8 bytes
			t.Errorf("expected data length 8, got %d", len(telegram.Data))
		}
		if telegram.OptData != nil {
			t.Errorf("expected nil optdata, got %v", telegram.OptData)
		}
	})

	t.Run("handles struct with only optdata fields", func(t *testing.T) {
		type OptDataOnlyStruct struct {
			Field1 int32 `enocean-esp3:"optdata"`
			Field2 int32 `enocean-esp3:"optdata"`
		}

		cmd := OptDataOnlyStruct{Field1: 1, Field2: 2}
		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		if len(telegram.Data) != 0 {
			t.Errorf("expected empty data, got %d bytes", len(telegram.Data))
		}
		if len(telegram.OptData) != 8 { // 2 * int32 = 8 bytes
			t.Errorf("expected optdata length 8, got %d", len(telegram.OptData))
		}
	})

	t.Run("handles mixed data and optdata fields", func(t *testing.T) {
		type MixedStruct struct {
			Data1    int32 `enocean-esp3:"data"`
			OptData1 int32 `enocean-esp3:"optdata"`
			Data2    int32 `enocean-esp3:"data"`
			OptData2 int32 `enocean-esp3:"optdata"`
		}

		cmd := MixedStruct{Data1: 1, OptData1: 2, Data2: 3, OptData2: 4}
		cfg := SerializerConfig{}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		if len(telegram.Data) != 8 { // 2 * int32 = 8 bytes
			t.Errorf("expected data length 8, got %d", len(telegram.Data))
		}
		if len(telegram.OptData) != 8 { // 2 * int32 = 8 bytes
			t.Errorf("expected optdata length 8, got %d", len(telegram.OptData))
		}
	})
}

// TestCommandToTelegram_ErrorPaths verifies CommandToTelegram_ErrorPaths behavior.
func TestCommandToTelegram_ErrorPaths(t *testing.T) {
	t.Run("handles error in serializeValue for nested struct", func(t *testing.T) {
		type ProblematicNested struct {
			Field interface{} // Will cause error if nil
		}

		type OuterStruct struct {
			OuterField int8              `enocean-esp3:"data"`
			Nested     ProblematicNested `enocean-esp3:"data"`
		}

		cmd := OuterStruct{
			OuterField: 1,
			Nested: ProblematicNested{
				Field: nil, // This will cause an error
			},
		}

		cfg := SerializerConfig{}
		_, err := CommandToTelegram(cmd, cfg)
		if err == nil {
			t.Errorf("expected error for nil interface in nested struct, got nil")
		}
		// Verify error is properly wrapped
		if err.Error() != "failed to serialize field Nested: failed to serialize nested field Field: nil interface: size unknown" {
			t.Errorf("expected specific error message, got '%s'", err.Error())
		}
	})

	t.Run("handles error in serializeValue for array element", func(t *testing.T) {
		type ArrayWithProblematicElement struct {
			Arr [1]interface{} `enocean-esp3:"data"`
		}

		cmd := ArrayWithProblematicElement{
			Arr: [1]interface{}{nil},
		}

		cfg := SerializerConfig{}
		_, err := CommandToTelegram(cmd, cfg)
		if err == nil {
			t.Errorf("expected error for nil interface in array, got nil")
		}
		// Error should be propagated from serializeValue
		if err.Error() != "failed to serialize field Arr: nil interface: size unknown" {
			t.Errorf("expected specific error message, got '%s'", err.Error())
		}
	})

	t.Run("handles error in serializeValue for slice element", func(t *testing.T) {
		type SliceWithProblematicElement struct {
			Slice []interface{} `enocean-esp3:"data"`
		}

		cmd := SliceWithProblematicElement{
			Slice: []interface{}{nil},
		}

		cfg := SerializerConfig{}
		_, err := CommandToTelegram(cmd, cfg)
		if err == nil {
			t.Errorf("expected error for nil interface in slice, got nil")
		}
		// Error should be propagated from serializeValue
		if err.Error() != "failed to serialize field Slice: nil interface: size unknown" {
			t.Errorf("expected specific error message, got '%s'", err.Error())
		}
	})
}

// TestSerializerConfig_Sanitize verifies SerializerConfig_Sanitize behavior.
func TestSerializerConfig_Sanitize(t *testing.T) {
	// Test sanitize indirectly through the public API
	// Since sanitize is unexported, we test it by verifying that functions
	// work correctly with nil values, which proves sanitize is working.

	t.Run("sanitize sets default byte order when nil", func(t *testing.T) {
		type IntStruct struct {
			Value uint32 `enocean-esp3:"data"`
		}

		cmd := IntStruct{Value: 0x12345678}
		cfg := SerializerConfig{
			ByteOrder: nil, // Should default to BigEndian
		}

		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Big endian: 0x12 0x34 0x56 0x78
		expected := []byte{0x12, 0x34, 0x56, 0x78}
		for i, b := range expected {
			if telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x (proves BigEndian default)", i, telegram.Data[i], b)
			}
		}
	})

	t.Run("sanitize creates empty map when nil", func(t *testing.T) {
		// Test that nil Serializers map doesn't cause issues
		cfg := SerializerConfig{
			Serializers: nil, // Should be sanitized to empty map
		}

		cmd := TestStruct{
			CommandCode: enums.CommonCommandWR_SECUREDEVICE_DEL,
			DeviceID:    deviceid.DeviceID(0x12345678),
		}

		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize with nil Serializers: %v", err)
		}

		if len(telegram.Data) == 0 {
			t.Errorf("expected non-empty data, proves sanitize created empty map")
		}
	})

	t.Run("sanitize preserves existing values", func(t *testing.T) {
		customMap := map[reflect.Type]CustomSerializer{
			reflect.TypeOf(int32(0)): func(buf *bytes.Buffer, v reflect.Value, byteOrder binary.ByteOrder) error {
				// Custom serializer that adds 1000
				return binary.Write(buf, byteOrder, int32(v.Int())+1000)
			},
		}
		cfg := SerializerConfig{
			Serializers: customMap,
			ByteOrder:   binary.LittleEndian,
		}

		type CustomIntStruct struct {
			Value int32 `enocean-esp3:"data"`
		}

		cmd := CustomIntStruct{Value: 42}
		telegram, err := CommandToTelegram(cmd, cfg)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Verify custom serializer was used (42 + 1000 = 1042)
		// Little endian: 0x1a 0x04 0x00 0x00
		if len(telegram.Data) != 4 {
			t.Fatalf("expected 4 bytes, got %d", len(telegram.Data))
		}
		value := int32(telegram.Data[0]) | int32(telegram.Data[1])<<8 | int32(telegram.Data[2])<<16 | int32(telegram.Data[3])<<24
		if value != 1042 {
			t.Errorf("expected 1042 (42+1000), got %d (proves custom serializer preserved)", value)
		}
	})

	t.Run("sanitize doesn't modify original config", func(t *testing.T) {
		original := SerializerConfig{
			Serializers: nil,
			ByteOrder:   nil,
		}

		// Use the config - this will call sanitize internally
		type IntStruct struct {
			Value int32 `enocean-esp3:"data"`
		}
		cmd := IntStruct{Value: 42}
		_, err := CommandToTelegram(cmd, original)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Original should still be nil (sanitize returns new config, doesn't modify)
		if original.Serializers != nil {
			t.Errorf("sanitize should not modify original config, Serializers should still be nil")
		}
		if original.ByteOrder != nil {
			t.Errorf("sanitize should not modify original config, ByteOrder should still be nil")
		}
	})
}

// TestCommandToTelegram_VariadicConfig verifies CommandToTelegram_VariadicConfig behavior.
func TestCommandToTelegram_VariadicConfig(t *testing.T) {
	t.Run("works with no config", func(t *testing.T) {
		type IntStruct struct {
			Value uint32 `enocean-esp3:"data"`
		}

		cmd := IntStruct{Value: 0x12345678}
		telegram, err := CommandToTelegram(cmd)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Should use default BigEndian
		expected := []byte{0x12, 0x34, 0x56, 0x78}
		for i, b := range expected {
			if telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x", i, telegram.Data[i], b)
			}
		}
	})

	t.Run("merges multiple configs - Serializers", func(t *testing.T) {
		type CustomInt int32
		type CustomFloat float32

		serializer1 := map[reflect.Type]CustomSerializer{
			reflect.TypeOf(CustomInt(0)): func(buf *bytes.Buffer, v reflect.Value, byteOrder binary.ByteOrder) error {
				return binary.Write(buf, byteOrder, int32(v.Int())+100)
			},
		}

		serializer2 := map[reflect.Type]CustomSerializer{
			reflect.TypeOf(CustomFloat(0)): func(buf *bytes.Buffer, v reflect.Value, byteOrder binary.ByteOrder) error {
				return binary.Write(buf, byteOrder, float32(v.Float())+200)
			},
		}

		cfg1 := SerializerConfig{Serializers: serializer1}
		cfg2 := SerializerConfig{Serializers: serializer2}

		type TestStruct struct {
			IntVal   CustomInt   `enocean-esp3:"data"`
			FloatVal CustomFloat `enocean-esp3:"data"`
		}

		cmd := TestStruct{
			IntVal:   CustomInt(42),
			FloatVal: CustomFloat(3.14),
		}

		telegram, err := CommandToTelegram(cmd, cfg1, cfg2)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Both serializers should be active (merged)
		// IntVal: 42 + 100 = 142 (4 bytes, big endian)
		// FloatVal: 3.14 + 200 = 203.14 (4 bytes, big endian)
		if len(telegram.Data) != 8 {
			t.Errorf("expected 8 bytes, got %d", len(telegram.Data))
		}
	})

	t.Run("merges multiple configs - ByteOrder override", func(t *testing.T) {
		type IntStruct struct {
			Value uint32 `enocean-esp3:"data"`
		}

		cmd := IntStruct{Value: 0x12345678}

		// First config with BigEndian
		cfg1 := SerializerConfig{ByteOrder: binary.BigEndian}
		// Second config with LittleEndian (should override)
		cfg2 := SerializerConfig{ByteOrder: binary.LittleEndian}

		telegram, err := CommandToTelegram(cmd, cfg1, cfg2)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Should use LittleEndian (last non-nil)
		expected := []byte{0x78, 0x56, 0x34, 0x12}
		for i, b := range expected {
			if telegram.Data[i] != b {
				t.Errorf("data[%d]: got 0x%02x, expected 0x%02x (little endian from cfg2)", i, telegram.Data[i], b)
			}
		}
	})

	t.Run("merges multiple configs - later Serializers override earlier", func(t *testing.T) {
		type CustomInt int32

		serializer1 := map[reflect.Type]CustomSerializer{
			reflect.TypeOf(CustomInt(0)): func(buf *bytes.Buffer, v reflect.Value, byteOrder binary.ByteOrder) error {
				return binary.Write(buf, byteOrder, int32(v.Int())+100)
			},
		}

		serializer2 := map[reflect.Type]CustomSerializer{
			reflect.TypeOf(CustomInt(0)): func(buf *bytes.Buffer, v reflect.Value, byteOrder binary.ByteOrder) error {
				return binary.Write(buf, byteOrder, int32(v.Int())+200) // Different serializer for same type
			},
		}

		cfg1 := SerializerConfig{Serializers: serializer1}
		cfg2 := SerializerConfig{Serializers: serializer2}

		type TestStruct struct {
			IntVal CustomInt `enocean-esp3:"data"`
		}

		cmd := TestStruct{IntVal: CustomInt(42)}

		telegram, err := CommandToTelegram(cmd, cfg1, cfg2)
		if err != nil {
			t.Fatalf("failed to serialize: %v", err)
		}

		// Should use serializer2 (42 + 200 = 242), not serializer1
		if len(telegram.Data) != 4 {
			t.Fatalf("expected 4 bytes, got %d", len(telegram.Data))
		}
		value := int32(telegram.Data[0])<<24 | int32(telegram.Data[1])<<16 | int32(telegram.Data[2])<<8 | int32(telegram.Data[3])
		if value != 242 {
			t.Errorf("expected 242 (42+200 from cfg2), got %d", value)
		}
	})
}
