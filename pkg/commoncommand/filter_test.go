package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

func TestNewWrFilterAdd(t *testing.T) {
	t.Run("creates filter add command with forward and repeat", func(t *testing.T) {
		cmd, err := NewWrFilterAdd(enums.FilterCriterionRSSI, 0x12345678, true, true)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_FILTER_ADD {
			t.Errorf("expected CommandCode WR_FILTER_ADD, got 0x%02x", cmd.CommandCode)
		}

		expectedAction := byte(enums.FilterActionFORWARD | enums.FilterActionREPEAT)
		if cmd.Action != expectedAction {
			t.Errorf("expected Action 0x%02x (forward+repeat), got 0x%02x", expectedAction, cmd.Action)
		}

		if cmd.Criterion != enums.FilterCriterionRSSI {
			t.Errorf("expected Criterion RSSI, got %v", cmd.Criterion)
		}

		if cmd.Value != 0x12345678 {
			t.Errorf("expected Value 0x%08x, got 0x%08x", 0x12345678, cmd.Value)
		}
	})

	t.Run("creates filter add command without forward and without repeat", func(t *testing.T) {
		cmd, err := NewWrFilterAdd(enums.FilterCriterionRORG, 0xAABBCCDD, false, false)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expectedAction := byte(enums.FilterActionNO_FORWARD | enums.FilterActionNO_REPEAT)
		if cmd.Action != expectedAction {
			t.Errorf("expected Action 0x%02x (no_forward+no_repeat), got 0x%02x", expectedAction, cmd.Action)
		}
	})

	t.Run("creates filter add command with mixed actions", func(t *testing.T) {
		cmd, err := NewWrFilterAdd(enums.FilterCriterionSENDER_ID, 0x00000001, true, false)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expectedAction := byte(enums.FilterActionFORWARD | enums.FilterActionNO_REPEAT)
		if cmd.Action != expectedAction {
			t.Errorf("expected Action 0x%02x (forward+no_repeat), got 0x%02x", expectedAction, cmd.Action)
		}
	})
}

func TestWrFilterAdd_Serialize(t *testing.T) {
	t.Run("serializes filter add command", func(t *testing.T) {
		cmd, _ := NewWrFilterAdd(enums.FilterCriterionRSSI, 0x12345678, true, true)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if telegram.PacketType != enums.PacketTypeCOMMON_COMMAND {
			t.Errorf("expected PacketType COMMON_COMMAND, got %v", telegram.PacketType)
		}

		// Data: Command(1) + Action(1) + Criterion(1) + Value(4) = 7 bytes
		if len(telegram.Data) != 7 {
			t.Errorf("expected Data length 7, got %d", len(telegram.Data))
		}

		// Check command code
		if telegram.Data[0] != byte(enums.CommonCommandWR_FILTER_ADD) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_FILTER_ADD, telegram.Data[0])
		}
	})
}

func TestNewWrFilterDel(t *testing.T) {
	t.Run("creates filter del command with forward and repeat", func(t *testing.T) {
		cmd, err := NewWrFilterDel(enums.FilterCriterionDESTINATION_ID, 0x87654321, true, true)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_FILTER_DEL {
			t.Errorf("expected CommandCode WR_FILTER_DEL, got 0x%02x", cmd.CommandCode)
		}

		expectedAction := byte(enums.FilterActionFORWARD | enums.FilterActionREPEAT)
		if cmd.Action != expectedAction {
			t.Errorf("expected Action 0x%02x, got 0x%02x", expectedAction, cmd.Action)
		}
	})

	t.Run("creates filter del command without forward and without repeat", func(t *testing.T) {
		cmd, err := NewWrFilterDel(enums.FilterCriterionRORG, 0xAABBCCDD, false, false)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expectedAction := byte(enums.FilterActionNO_FORWARD | enums.FilterActionNO_REPEAT)
		if cmd.Action != expectedAction {
			t.Errorf("expected Action 0x%02x, got 0x%02x", expectedAction, cmd.Action)
		}
	})
}

func TestWrFilterDel_Serialize(t *testing.T) {
	t.Run("serializes filter del command", func(t *testing.T) {
		cmd, _ := NewWrFilterDel(enums.FilterCriterionRORG, 0xAABBCCDD, false, false)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Action(1) + Criterion(1) + Value(4) = 7 bytes
		if len(telegram.Data) != 7 {
			t.Errorf("expected Data length 7, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_FILTER_DEL) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_FILTER_DEL, telegram.Data[0])
		}
	})
}

func TestNewWrFilterDelAll(t *testing.T) {
	t.Run("creates filter delete all command", func(t *testing.T) {
		cmd, err := NewWrFilterDelAll()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_FILTER_DEL_ALL {
			t.Errorf("expected CommandCode WR_FILTER_DEL_ALL, got 0x%02x", cmd.CommandCode)
		}
	})
}

func TestWrFilterDelAll_Serialize(t *testing.T) {
	t.Run("serializes filter del all command", func(t *testing.T) {
		cmd, _ := NewWrFilterDelAll()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) = 1 byte
		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_FILTER_DEL_ALL) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_FILTER_DEL_ALL, telegram.Data[0])
		}
	})
}

func TestNewWrFilterEnable(t *testing.T) {
	t.Run("creates filter enable command with OR_ALL_FILTERS operator", func(t *testing.T) {
		cmd, err := NewWrFilterEnable(true, enums.FilerOperatorOR_ALL_FILTERS)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_FILTER_ENABLE {
			t.Errorf("expected CommandCode WR_FILTER_ENABLE, got 0x%02x", cmd.CommandCode)
		}

		if !cmd.Toggle {
			t.Errorf("expected Toggle = true, got false")
		}

		if cmd.FilerOperator != enums.FilerOperatorOR_ALL_FILTERS {
			t.Errorf("expected FilerOperator OR_ALL_FILTERS, got %v", cmd.FilerOperator)
		}
	})

	t.Run("creates filter enable command with AND_ALL_FILTERS operator", func(t *testing.T) {
		cmd, err := NewWrFilterEnable(false, enums.FilerOperatorAND_ALL_FILTERS)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.Toggle {
			t.Errorf("expected Toggle = false, got true")
		}

		if cmd.FilerOperator != enums.FilerOperatorAND_ALL_FILTERS {
			t.Errorf("expected FilerOperator AND_ALL_FILTERS, got %v", cmd.FilerOperator)
		}
	})
}

func TestWrFilterEnable_Serialize(t *testing.T) {
	t.Run("serializes filter enable command", func(t *testing.T) {
		cmd, _ := NewWrFilterEnable(true, enums.FilerOperatorOR_ALL_FILTERS)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Toggle(1) + Operator(1) = 3 bytes
		if len(telegram.Data) != 3 {
			t.Errorf("expected Data length 3, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandWR_FILTER_ENABLE) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandWR_FILTER_ENABLE, telegram.Data[0])
		}

		if telegram.Data[1] != 0x01 {
			t.Errorf("expected Data[1] (Toggle) = 0x01, got 0x%02x", telegram.Data[1])
		}
	})
}

func TestNewRdFilter(t *testing.T) {
	t.Run("creates read filter command", func(t *testing.T) {
		cmd, err := NewRdFilter()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_FILTER {
			t.Errorf("expected CommandCode RD_FILTER, got 0x%02x", cmd.CommandCode)
		}
	})
}

func TestRdFilter_Serialize(t *testing.T) {
	t.Run("serializes read filter command", func(t *testing.T) {
		cmd, _ := NewRdFilter()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) = 1 byte
		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandRD_FILTER) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandRD_FILTER, telegram.Data[0])
		}
	})
}

func TestParseRdFilterResponseOK(t *testing.T) {
	t.Run("parses filter response with no filters", func(t *testing.T) {
		// Count = 0, no filter data
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x00},
			OptData: nil,
		}

		result, err := ParseRdFilterResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(result.Filters) != 0 {
			t.Errorf("expected 0 filters, got %d", len(result.Filters))
		}
	})

	t.Run("parses filter response with multiple filters", func(t *testing.T) {
		resp := response.Packet{Code: enums.ReturnCodeSUCCESS, Data: []byte{
			0x02,
			0x00, 0x00, 0x00, 0x00, 0x01,
			0x01, 0x00, 0x00, 0x00, 0x02,
		}}
		result, err := ParseRdFilterResponseOK(resp)
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Filters) != 2 || result.Filters[0].Value != 1 || result.Filters[1].Value != 2 {
			t.Fatalf("filters = %#v", result.Filters)
		}
	})

	t.Run("rejects filter count mismatches", func(t *testing.T) {
		for _, data := range [][]byte{
			{0x02, 0x00, 0, 0, 0, 1},
			{0x02, 0x00, 0, 0, 0, 1, 0x01, 0, 0, 0, 2, 0x02, 0, 0, 0, 3},
		} {
			if _, err := ParseRdFilterResponseOK(response.Packet{Code: enums.ReturnCodeSUCCESS, Data: data}); err == nil {
				t.Fatalf("count mismatch accepted: %x", data)
			}
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x00},
			OptData: nil,
		}

		_, err := ParseRdFilterResponseOK(resp)
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

		_, err := ParseRdFilterResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}
	})
}
