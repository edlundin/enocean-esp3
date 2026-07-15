package commoncommand

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"testing"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

// TestNewRdVersion verifies NewRdVersion behavior.
func TestNewRdVersion(t *testing.T) {
	t.Run("creates read version command", func(t *testing.T) {
		cmd, err := NewRdVersion()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_VERSION {
			t.Errorf("expected CommandCode RD_VERSION, got 0x%02x", cmd.CommandCode)
		}
	})
}

// TestRdVersion_Serialize verifies RdVersion_Serialize behavior.
func TestRdVersion_Serialize(t *testing.T) {
	t.Run("serializes read version command", func(t *testing.T) {
		cmd, err := NewRdVersion()
		if err != nil {
			t.Fatalf("expected no constructor error, got: %v", err)
		}
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Fatalf("expected Data length 1, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandRD_VERSION) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandRD_VERSION, telegram.Data[0])
		}
	})
}

// TestParseRdVersionResponseOK verifies ParseRdVersionResponseOK behavior.
func TestParseRdVersionResponseOK(t *testing.T) {
	t.Run("parses version response", func(t *testing.T) {
		// Response: AppVersion(4) + ApiVersion(4) + ChipID(4) + ChipVersion(4) + Description(variable)
		// Test with a minimal description
		resp := response.Packet{
			Code: enums.ReturnCodeSUCCESS,
			Data: []byte{
				0x01, 0x02, 0x03, 0x04, // AppVersion
				0x05, 0x06, 0x07, 0x08, // ApiVersion
				0x12, 0x34, 0x56, 0x78, // ChipID (big-endian uint32: 0x12345678)
				0xAA, 0xBB, 0xCC, 0xDD, // ChipVersion
				'T', 'C', 'M', '3', '1', '5', // Description
			},
			OptData: nil,
		}

		result, err := ParseRdVersionResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expectedAppVersion := [4]byte{0x01, 0x02, 0x03, 0x04}
		if result.AppVersion != expectedAppVersion {
			t.Errorf("expected AppVersion %v, got %v", expectedAppVersion, result.AppVersion)
		}

		expectedApiVersion := [4]byte{0x05, 0x06, 0x07, 0x08}
		if result.ApiVersion != expectedApiVersion {
			t.Errorf("expected ApiVersion %v, got %v", expectedApiVersion, result.ApiVersion)
		}

		if result.ChipID != 0x12345678 {
			t.Errorf("expected ChipID 0x12345678, got 0x%08x", result.ChipID)
		}

		if result.ChipVersion != 0xAABBCCDD {
			t.Errorf("expected ChipVersion 0xAABBCCDD, got 0x%08x", result.ChipVersion)
		}

		if result.Description != "TCM315" {
			t.Errorf("expected Description 'TCM315', got '%s'", result.Description)
		}
	})

	t.Run("parses version response with empty description", func(t *testing.T) {
		resp := response.Packet{
			Code: enums.ReturnCodeSUCCESS,
			Data: []byte{
				0x01, 0x02, 0x03, 0x04, // AppVersion
				0x05, 0x06, 0x07, 0x08, // ApiVersion
				0x00, 0x00, 0x00, 0x00, // ChipID
				0x00, 0x00, 0x00, 0x00, // ChipVersion
				// Empty description
			},
			OptData: nil,
		}

		result, err := ParseRdVersionResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.Description != "" {
			t.Errorf("expected empty Description, got '%s'", result.Description)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01, 0x02, 0x03, 0x04},
			OptData: nil,
		}

		_, err := ParseRdVersionResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseRdVersionResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}
	})

	t.Run("uses custom deserializer for string field", func(t *testing.T) {
		// Test that the custom string deserializer reads all remaining bytes
		resp := response.Packet{
			Code: enums.ReturnCodeSUCCESS,
			Data: []byte{
				0x01, 0x02, 0x03, 0x04, // AppVersion
				0x05, 0x06, 0x07, 0x08, // ApiVersion
				0x12, 0x34, 0x56, 0x78, // ChipID
				0xAA, 0xBB, 0xCC, 0xDD, // ChipVersion
				'V', 'e', 'r', 's', 'i', 'o', 'n', ':', ' ', '1', '.', '0', // Description
			},
			OptData: nil,
		}

		result, err := ParseRdVersionResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.Description != "Version: 1.0" {
			t.Errorf("expected Description 'Version: 1.0', got '%s'", result.Description)
		}
	})
}

// TestVersionStringDeserializer verifies custom string deserialization.
func TestVersionStringDeserializer(t *testing.T) {
	t.Run("custom string deserializer reads all remaining bytes", func(t *testing.T) {
		stringDeserializer := func(buf *bytes.Reader, v reflect.Value, _ binary.ByteOrder) error {
			rest, err := io.ReadAll(buf)
			if err != nil {
				return err
			}
			v.SetString(string(rest))
			return nil
		}

		cfg := serializer.DeserializerConfig{
			Deserializers: map[reflect.Type]serializer.CustomDeserializer{
				reflect.TypeOf(""): stringDeserializer,
			},
		}

		// Test struct with a string field
		type TestStruct struct {
			Value   uint32
			Message string
		}

		data := []byte{
			0x00, 0x00, 0x00, 0x01, // Value = 1
			'H', 'e', 'l', 'l', 'o', // Message = "Hello"
		}

		var result TestStruct
		err := serializer.BytesToStruct(data, &result, cfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.Value != 1 {
			t.Errorf("expected Value 1, got %d", result.Value)
		}

		if result.Message != "Hello" {
			t.Errorf("expected Message 'Hello', got '%s'", result.Message)
		}
	})
}
