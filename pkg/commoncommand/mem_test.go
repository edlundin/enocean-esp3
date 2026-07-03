package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

func TestNewWrMem(t *testing.T) {
	t.Run("creates write memory command", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03, 0x04}
		cmd, err := NewWrMem(enums.MemoryTypeFLASH, 0x00100000, data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_MEM {
			t.Errorf("expected CommandCode WR_MEM, got 0x%02x", cmd.CommandCode)
		}

		if cmd.Type != enums.MemoryTypeFLASH {
			t.Errorf("expected Type FLASH, got %v", cmd.Type)
		}

		if cmd.Address != 0x00100000 {
			t.Errorf("expected Address 0x%08x, got 0x%08x", 0x00100000, cmd.Address)
		}

		if len(cmd.Data) != 4 {
			t.Errorf("expected Data length 4, got %d", len(cmd.Data))
		}
	})

	t.Run("creates write memory command with RAM type", func(t *testing.T) {
		data := []byte{0xAA, 0xBB}
		cmd, err := NewWrMem(enums.MemoryTypeRAM0, 0x00200000, data)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.Type != enums.MemoryTypeRAM0 {
			t.Errorf("expected Type RAM0, got %v", cmd.Type)
		}
	})
}

func TestWrMem_Serialize(t *testing.T) {
	t.Run("serializes write memory command", func(t *testing.T) {
		data := []byte{0x01, 0x02, 0x03}
		cmd, _ := NewWrMem(enums.MemoryTypeFLASH, 0x00100000, data)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Type(1) + Address(4) + Data(variable)
		// Minimum: 1 + 1 + 4 = 6 bytes
		if len(telegram.Data) < 6 {
			t.Errorf("expected Data length >= 6, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_MEM) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_MEM, telegram.Data[0])
		}
	})
}

func TestNewRdMem(t *testing.T) {
	t.Run("creates read memory command", func(t *testing.T) {
		cmd, err := NewRdMem(enums.MemoryTypeFLASH, 0x00100000, 256)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_MEM {
			t.Errorf("expected CommandCode RD_MEM, got 0x%02x", cmd.CommandCode)
		}

		if cmd.Type != enums.MemoryTypeFLASH {
			t.Errorf("expected Type FLASH, got %v", cmd.Type)
		}

		if cmd.Address != 0x00100000 {
			t.Errorf("expected Address 0x%08x, got 0x%08x", 0x00100000, cmd.Address)
		}

		if cmd.DataLength != 256 {
			t.Errorf("expected DataLength 256, got %d", cmd.DataLength)
		}
	})
}

func TestRdMem_Serialize(t *testing.T) {
	t.Run("serializes read memory command", func(t *testing.T) {
		cmd, _ := NewRdMem(enums.MemoryTypeRAM0, 0x00200000, 128)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Type(1) + Address(4) + DataLength(2) = 8 bytes
		if len(telegram.Data) != 8 {
			t.Errorf("expected Data length 8, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandRD_MEM) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandRD_MEM, telegram.Data[0])
		}
	})
}

func TestParseRdMemResponseOK(t *testing.T) {
	t.Run("parses read memory response", func(t *testing.T) {
		// Response: Data(variable bytes)
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			OptData: nil,
		}

		result, err := ParseRdMemResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result.Data) != 5 {
			t.Errorf("expected Data length 5, got %d", len(result.Data))
		}

		expected := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
		for i, b := range expected {
			if result.Data[i] != b {
				t.Errorf("Data[%d]: expected 0x%02x, got 0x%02x", i, b, result.Data[i])
			}
		}
	})

	t.Run("parses empty memory response", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		result, err := ParseRdMemResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error for empty data, got: %v", err)
		}

		if len(result.Data) != 0 {
			t.Errorf("expected empty Data, got %d bytes", len(result.Data))
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01, 0x02, 0x03},
			OptData: nil,
		}

		_, err := ParseRdMemResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})
}

func TestNewRdMemAddress(t *testing.T) {
	t.Run("creates read memory address command", func(t *testing.T) {
		cmd, err := NewRdMemAddress(enums.MemoryAreaCONFIG)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_MEM_ADDRESS {
			t.Errorf("expected CommandCode RD_MEM_ADDRESS, got 0x%02x", cmd.CommandCode)
		}

		if cmd.Area != enums.MemoryAreaCONFIG {
			t.Errorf("expected Area CONFIG, got %v", cmd.Area)
		}
	})

	t.Run("creates read memory address command with SMART_ACK_TABLE area", func(t *testing.T) {
		cmd, err := NewRdMemAddress(enums.MemoryAreaSMART_ACK_TABLE)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.Area != enums.MemoryAreaSMART_ACK_TABLE {
			t.Errorf("expected Area SMART_ACK_TABLE, got %v", cmd.Area)
		}
	})
}

func TestRdMemAddress_Serialize(t *testing.T) {
	t.Run("serializes read memory address command", func(t *testing.T) {
		cmd, _ := NewRdMemAddress(enums.MemoryAreaCONFIG)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Area(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandRD_MEM_ADDRESS) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandRD_MEM_ADDRESS, telegram.Data[0])
		}
	})
}

func TestParseRdMemAddressResponseOK(t *testing.T) {
	t.Run("parses read memory address response", func(t *testing.T) {
		// Response: Type(1) + Address(4) + Length(4) = 9 bytes
		// Big-endian format: high byte first
		resp := response.Packet{
			Code: enums.ReturnCodeSUCCESS,
			Data: []byte{
				0x00,                   // Type = FLASH
				0x00, 0x10, 0x00, 0x00, // Address = 0x00100000 (big-endian)
				0x00, 0x01, 0x00, 0x00, // Length = 0x00010000 (big-endian)
			},
			OptData: nil,
		}

		result, err := ParseRdMemAddressResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.Type != enums.MemoryTypeFLASH {
			t.Errorf("expected Type FLASH, got %v", result.Type)
		}

		if result.Address != 0x00100000 {
			t.Errorf("expected Address 0x%08x, got 0x%08x", 0x00100000, result.Address)
		}

		if result.Length != 0x00010000 {
			t.Errorf("expected Length 0x%08x, got 0x%08x", 0x00010000, result.Length)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00},
			OptData: nil,
		}

		_, err := ParseRdMemAddressResponseOK(resp)
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

		_, err := ParseRdMemAddressResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})

	t.Run("returns error for invalid memory type", func(t *testing.T) {
		resp := response.Packet{
			Code: enums.ReturnCodeSUCCESS,
			Data: []byte{
				0xFF,                   // Type = invalid
				0x00, 0x10, 0x00, 0x00, // Address
				0x00, 0x01, 0x00, 0x00, // Length
			},
			OptData: nil,
		}

		_, err := ParseRdMemAddressResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for invalid memory type, got nil")
		}

		if err.Error() != "invalid memory type" {
			t.Errorf("expected error 'invalid memory type', got '%s'", err.Error())
		}
	})
}
