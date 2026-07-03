package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

func TestNewRdSysLog(t *testing.T) {
	t.Run("creates read sys log command", func(t *testing.T) {
		cmd, err := NewRdSysLog()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_SYS_LOG {
			t.Errorf("expected CommandCode RD_SYS_LOG, got 0x%02x", cmd.CommandCode)
		}
	})
}

func TestRdSysLog_Serialize(t *testing.T) {
	t.Run("serializes read sys log command", func(t *testing.T) {
		cmd, _ := NewRdSysLog()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandRD_SYS_LOG) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandRD_SYS_LOG, telegram.Data[0])
		}
	})
}

func TestParseRdSysLogResponseOK(t *testing.T) {
	t.Run("parses sys log response", func(t *testing.T) {
		// Response: ApiLogEntries in Data, AppLogEntries in OptData
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			OptData: []byte{0x06, 0x07, 0x08},
		}

		result, err := ParseRdSysLogResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result.ApiLogEntries) != 5 {
			t.Errorf("expected ApiLogEntries length 5, got %d", len(result.ApiLogEntries))
		}

		if len(result.AppLogEntries) != 3 {
			t.Errorf("expected AppLogEntries length 3, got %d", len(result.AppLogEntries))
		}

		// Verify data content
		expectedApiLog := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
		for i, b := range expectedApiLog {
			if result.ApiLogEntries[i] != b {
				t.Errorf("ApiLogEntries[%d]: expected 0x%02x, got 0x%02x", i, b, result.ApiLogEntries[i])
			}
		}

		expectedAppLog := []byte{0x06, 0x07, 0x08}
		for i, b := range expectedAppLog {
			if result.AppLogEntries[i] != b {
				t.Errorf("AppLogEntries[%d]: expected 0x%02x, got 0x%02x", i, b, result.AppLogEntries[i])
			}
		}
	})

	t.Run("parses sys log response with empty data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: []byte{},
		}

		result, err := ParseRdSysLogResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result.ApiLogEntries) != 0 {
			t.Errorf("expected empty ApiLogEntries, got %d bytes", len(result.ApiLogEntries))
		}

		if len(result.AppLogEntries) != 0 {
			t.Errorf("expected empty AppLogEntries, got %d bytes", len(result.AppLogEntries))
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01, 0x02},
			OptData: []byte{0x03},
		}

		_, err := ParseRdSysLogResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}

		if err.Error() != "invalid return code" {
			t.Errorf("expected error 'invalid return code', got '%s'", err.Error())
		}
	})
}

func TestNewResetSysLog(t *testing.T) {
	t.Run("creates reset sys log command", func(t *testing.T) {
		cmd, err := NewResetSysLog()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRESET_SYS_LOG {
			t.Errorf("expected CommandCode RESET_SYS_LOG, got 0x%02x", cmd.CommandCode)
		}
	})
}

func TestResetSysLog_Serialize(t *testing.T) {
	t.Run("serializes reset sys log command", func(t *testing.T) {
		cmd, _ := NewResetSysLog()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandRESET_SYS_LOG) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandRESET_SYS_LOG, telegram.Data[0])
		}
	})
}
