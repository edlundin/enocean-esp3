package commoncommand

import (
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

// TestNewRdDutyCycleLimit verifies NewRdDutyCycleLimit behavior.
func TestNewRdDutyCycleLimit(t *testing.T) {
	t.Run("creates read duty cycle limit command", func(t *testing.T) {
		cmd, err := NewRdDutyCycleLimit()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_DUTYCYCLE_LIMIT {
			t.Errorf("expected CommandCode RD_DUTYCYCLE_LIMIT, got 0x%02x", cmd.CommandCode)
		}
	})
}

// TestRdDutyCycleLimit_Serialize verifies RdDutyCycleLimit_Serialize behavior.
func TestRdDutyCycleLimit_Serialize(t *testing.T) {
	t.Run("serializes read duty cycle limit command", func(t *testing.T) {
		cmd, _ := NewRdDutyCycleLimit()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandRD_DUTYCYCLE_LIMIT) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandRD_DUTYCYCLE_LIMIT, telegram.Data[0])
		}
	})
}

// TestParseRdDutyCycleLimitResponseOK verifies ParseRdDutyCycleLimitResponseOK behavior.
func TestParseRdDutyCycleLimitResponseOK(t *testing.T) {
	t.Run("parses duty cycle limit response", func(t *testing.T) {
		// Response: AvailableDutyCycle(1) + Slots(1) + SlotPeriod(2) + TimeLeftInCurrentSlot(2) + AvailableDutyCycleAfterCurrentSlot(1) = 7 bytes
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x50, 0x04, 0x00, 0x64, 0x00, 0x32, 0x40},
			OptData: nil,
		}

		result, err := ParseRdDutyCycleLimitResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.AvailableDutyCycle != 0x50 {
			t.Errorf("expected AvailableDutyCycle = 0x50, got 0x%02x", result.AvailableDutyCycle)
		}

		if result.Slots != 0x04 {
			t.Errorf("expected Slots = 4, got %d", result.Slots)
		}

		if result.SlotPeriod != 0x0064 {
			t.Errorf("expected SlotPeriod = 100, got %d", result.SlotPeriod)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x50, 0x04, 0x00, 0x64, 0x00, 0x32, 0x40},
			OptData: nil,
		}

		_, err := ParseRdDutyCycleLimitResponseOK(resp)
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

		_, err := ParseRdDutyCycleLimitResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewSetBaudrate verifies NewSetBaudrate behavior.
func TestNewSetBaudrate(t *testing.T) {
	t.Run("creates set baudrate command", func(t *testing.T) {
		cmd, err := NewSetBaudrate(enums.TCMBaudrate115200)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandSET_BAUDRATE {
			t.Errorf("expected CommandCode SET_BAUDRATE, got 0x%02x", cmd.CommandCode)
		}

		if cmd.Baudrate != enums.TCMBaudrate115200 {
			t.Errorf("expected Baudrate 115200, got %v", cmd.Baudrate)
		}
	})

	t.Run("creates set baudrate command with 57600", func(t *testing.T) {
		cmd, err := NewSetBaudrate(enums.TCMBaudrate57600)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.Baudrate != enums.TCMBaudrate57600 {
			t.Errorf("expected Baudrate 57600, got %v", cmd.Baudrate)
		}
	})
}

// TestSetBaudrate_Serialize verifies SetBaudrate_Serialize behavior.
func TestSetBaudrate_Serialize(t *testing.T) {
	t.Run("serializes set baudrate command", func(t *testing.T) {
		cmd, _ := NewSetBaudrate(enums.TCMBaudrate115200)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + Baudrate(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}

		if telegram.Data[0] != byte(enums.CommonCommandSET_BAUDRATE) {
			t.Errorf("expected Data[0] = 0x%02x, got 0x%02x", enums.CommonCommandSET_BAUDRATE, telegram.Data[0])
		}
	})
}

// TestNewGetFrequencyInfo verifies NewGetFrequencyInfo behavior.
func TestNewGetFrequencyInfo(t *testing.T) {
	t.Run("creates get frequency info command", func(t *testing.T) {
		cmd, err := NewGetFrequencyInfo()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandGET_FREQUENCY_INFO {
			t.Errorf("expected CommandCode GET_FREQUENCY_INFO, got 0x%02x", cmd.CommandCode)
		}
	})
}

// TestGetFrequencyInfo_Serialize verifies GetFrequencyInfo_Serialize behavior.
func TestGetFrequencyInfo_Serialize(t *testing.T) {
	t.Run("serializes get frequency info command", func(t *testing.T) {
		cmd, _ := NewGetFrequencyInfo()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}
	})
}

// TestParseGetFrequencyInfoResponseOK verifies ParseGetFrequencyInfoResponseOK behavior.
func TestParseGetFrequencyInfoResponseOK(t *testing.T) {
	t.Run("parses frequency info response", func(t *testing.T) {
		// Response: Frequency(1) + Protocol(1) = 2 bytes
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01, 0x01},
			OptData: nil,
		}

		result, err := ParseGetFrequencyInfoResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.Frequency != enums.TCMFrequency868_000_MHZ {
			t.Errorf("expected Frequency 868.000MHz, got %v", result.Frequency)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01, 0x01},
			OptData: nil,
		}

		_, err := ParseGetFrequencyInfoResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseGetFrequencyInfoResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewGetStepCode verifies NewGetStepCode behavior.
func TestNewGetStepCode(t *testing.T) {
	t.Run("creates get stepcode command", func(t *testing.T) {
		cmd, err := NewGetStepCode()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandGET_STEPCODE {
			t.Errorf("expected CommandCode GET_STEPCODE, got 0x%02x", cmd.CommandCode)
		}
	})
}

// TestGetStepCode_Serialize verifies GetStepCode_Serialize behavior.
func TestGetStepCode_Serialize(t *testing.T) {
	t.Run("serializes get stepcode command", func(t *testing.T) {
		cmd, _ := NewGetStepCode()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}
	})
}

// TestParseGetStepCodeResponseOK verifies ParseGetStepCodeResponseOK behavior.
func TestParseGetStepCodeResponseOK(t *testing.T) {
	t.Run("parses stepcode response", func(t *testing.T) {
		// Response: StepCode(1) + Revision(1) = 2 bytes
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x02, 0x01},
			OptData: nil,
		}

		result, err := ParseGetStepCodeResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.StepCode != 0x02 {
			t.Errorf("expected StepCode = 2, got %d", result.StepCode)
		}

		if result.Revision != 0x01 {
			t.Errorf("expected Revision = 1, got %d", result.Revision)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x02, 0x01},
			OptData: nil,
		}

		_, err := ParseGetStepCodeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseGetStepCodeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewWrStartupDelay verifies NewWrStartupDelay behavior.
func TestNewWrStartupDelay(t *testing.T) {
	t.Run("creates write startup delay command", func(t *testing.T) {
		cmd, err := NewWrStartupDelay(50)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_STARTUP_DELAY {
			t.Errorf("expected CommandCode WR_STARTUP_DELAY, got 0x%02x", cmd.CommandCode)
		}

		if cmd.StartupDelay != 50 {
			t.Errorf("expected StartupDelay = 50, got %d", cmd.StartupDelay)
		}
	})
}

// TestWrStartupDelay_Serialize verifies WrStartupDelay_Serialize behavior.
func TestWrStartupDelay_Serialize(t *testing.T) {
	t.Run("serializes write startup delay command", func(t *testing.T) {
		cmd, _ := NewWrStartupDelay(100)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + StartupDelay(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}

		if telegram.Data[1] != 100 {
			t.Errorf("expected Data[1] = 100, got %d", telegram.Data[1])
		}
	})
}

// TestNewSetNoiseThreshold verifies NewSetNoiseThreshold behavior.
func TestNewSetNoiseThreshold(t *testing.T) {
	t.Run("creates set noise threshold command", func(t *testing.T) {
		cmd, err := NewSetNoiseThreshold(80)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandSET_NOISETHRESHOLD {
			t.Errorf("expected CommandCode SET_NOISETHRESHOLD, got 0x%02x", cmd.CommandCode)
		}

		if cmd.NoiseThreshold != 80 {
			t.Errorf("expected NoiseThreshold = 80, got %d", cmd.NoiseThreshold)
		}
	})
}

// TestSetNoiseThreshold_Serialize verifies SetNoiseThreshold_Serialize behavior.
func TestSetNoiseThreshold_Serialize(t *testing.T) {
	t.Run("serializes set noise threshold command", func(t *testing.T) {
		cmd, _ := NewSetNoiseThreshold(80)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + NoiseThreshold(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}
	})
}

// TestNewGetNoiseThreshold verifies NewGetNoiseThreshold behavior.
func TestNewGetNoiseThreshold(t *testing.T) {
	t.Run("creates get noise threshold command", func(t *testing.T) {
		cmd, err := NewGetNoiseThreshold()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandGET_NOISETHRESHOLD {
			t.Errorf("expected CommandCode GET_NOISETHRESHOLD, got 0x%02x", cmd.CommandCode)
		}
	})
}

// TestGetNoiseThreshold_Serialize verifies GetNoiseThreshold_Serialize behavior.
func TestGetNoiseThreshold_Serialize(t *testing.T) {
	t.Run("serializes get noise threshold command", func(t *testing.T) {
		cmd, _ := NewGetNoiseThreshold()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}
	})
}

// TestParseGetNoiseThresholdResponseOK verifies ParseGetNoiseThresholdResponseOK behavior.
func TestParseGetNoiseThresholdResponseOK(t *testing.T) {
	t.Run("parses noise threshold response", func(t *testing.T) {
		// Response: RSSILevel(1) = 1 byte
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x60},
			OptData: nil,
		}

		result, err := ParseGetNoiseThresholdResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.RSSILevel != 0x60 {
			t.Errorf("expected RSSILevel = 0x60, got 0x%02x", result.RSSILevel)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x60},
			OptData: nil,
		}

		_, err := ParseGetNoiseThresholdResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseGetNoiseThresholdResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewSetCRCMode verifies NewSetCRCMode behavior.
func TestNewSetCRCMode(t *testing.T) {
	t.Run("creates set CRC mode command", func(t *testing.T) {
		cmd, err := NewSetCRCMode(enums.CRCMode8BIT)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandSET_CRCMode {
			t.Errorf("expected CommandCode SET_CRCMode, got 0x%02x", cmd.CommandCode)
		}

		if cmd.CRCMode != enums.CRCMode8BIT {
			t.Errorf("expected CRCMode 8BIT, got %v", cmd.CRCMode)
		}
	})
}

// TestSetCRCMode_Serialize verifies SetCRCMode_Serialize behavior.
func TestSetCRCMode_Serialize(t *testing.T) {
	t.Run("serializes set CRC mode command", func(t *testing.T) {
		cmd, _ := NewSetCRCMode(enums.CRCMode8BIT)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + CRCMode(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}
	})
}

// TestNewGetCRCMode verifies NewGetCRCMode behavior.
func TestNewGetCRCMode(t *testing.T) {
	t.Run("creates get CRC mode command", func(t *testing.T) {
		cmd, err := NewGetCRCMode()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandGET_CRCMode {
			t.Errorf("expected CommandCode GET_CRCMode, got 0x%02x", cmd.CommandCode)
		}
	})
}

// TestGetCRCMode_Serialize verifies GetCRCMode_Serialize behavior.
func TestGetCRCMode_Serialize(t *testing.T) {
	t.Run("serializes get CRC mode command", func(t *testing.T) {
		cmd, _ := NewGetCRCMode()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}
	})
}

// TestParseGetCRCModeResponseOK verifies ParseGetCRCModeResponseOK behavior.
func TestParseGetCRCModeResponseOK(t *testing.T) {
	t.Run("parses CRC mode response", func(t *testing.T) {
		// Response: CRCMode(1) = 1 byte
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x00},
			OptData: nil,
		}

		result, err := ParseGetCRCModeResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.CRCMode != enums.CRCMode8BIT {
			t.Errorf("expected CRCMode 8BIT, got %v", result.CRCMode)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01},
			OptData: nil,
		}

		_, err := ParseGetCRCModeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{}, // Empty data - not enough for CRCMode(1)
			OptData: nil,
		}

		_, err := ParseGetCRCModeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewWrRSSITestMode verifies NewWrRSSITestMode behavior.
func TestNewWrRSSITestMode(t *testing.T) {
	t.Run("creates write RSSI test mode command", func(t *testing.T) {
		cmd, err := NewWrRSSITestMode(enums.RSSITestModeENABLED, 1000)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_RSSITEST_MODE {
			t.Errorf("expected CommandCode WR_RSSITEST_MODE, got 0x%02x", cmd.CommandCode)
		}

		if cmd.TestMode != enums.RSSITestModeENABLED {
			t.Errorf("expected TestMode ENABLED, got %v", cmd.TestMode)
		}

		if cmd.Timeout != 1000 {
			t.Errorf("expected Timeout = 1000, got %d", cmd.Timeout)
		}
	})
}

// TestWrRSSITestMode_Serialize verifies WrRSSITestMode_Serialize behavior.
func TestWrRSSITestMode_Serialize(t *testing.T) {
	t.Run("serializes write RSSI test mode command", func(t *testing.T) {
		cmd, _ := NewWrRSSITestMode(enums.RSSITestModeENABLED, 500)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + TestMode(1) + Timeout(2) = 4 bytes
		if len(telegram.Data) != 4 {
			t.Errorf("expected Data length 4, got %d", len(telegram.Data))
		}
	})
}

// TestNewRdRSSITestMode verifies NewRdRSSITestMode behavior.
func TestNewRdRSSITestMode(t *testing.T) {
	t.Run("creates read RSSI test mode command", func(t *testing.T) {
		cmd, err := NewRdRSSITestMode()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_RSSITEST_MODE {
			t.Errorf("expected CommandCode RD_RSSITEST_MODE, got 0x%02x", cmd.CommandCode)
		}
	})
}

// TestRdRSSITestMode_Serialize verifies RdRSSITestMode_Serialize behavior.
func TestRdRSSITestMode_Serialize(t *testing.T) {
	t.Run("serializes read RSSI test mode command", func(t *testing.T) {
		cmd, _ := NewRdRSSITestMode()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}
	})
}

// TestParseRdRSSITestModeResponseOK verifies ParseRdRSSITestModeResponseOK behavior.
func TestParseRdRSSITestModeResponseOK(t *testing.T) {
	t.Run("parses RSSI test mode response", func(t *testing.T) {
		// Response: TestMode(1) = 1 byte
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01},
			OptData: nil,
		}

		result, err := ParseRdRSSITestModeResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.TestMode != enums.RSSITestModeENABLED {
			t.Errorf("expected TestMode ENABLED, got %v", result.TestMode)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01},
			OptData: nil,
		}

		_, err := ParseRdRSSITestModeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseRdRSSITestModeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewWrTransparentMode verifies NewWrTransparentMode behavior.
func TestNewWrTransparentMode(t *testing.T) {
	t.Run("creates write transparent mode command", func(t *testing.T) {
		cmd, err := NewWrTransparentMode(enums.TransparentModeENABLED)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_TRANSPARENT_MODE {
			t.Errorf("expected CommandCode WR_TRANSPARENT_MODE, got 0x%02x", cmd.CommandCode)
		}

		if cmd.TransparentMode != enums.TransparentModeENABLED {
			t.Errorf("expected TransparentMode ENABLED, got %v", cmd.TransparentMode)
		}
	})
}

// TestWrTransparentMode_Serialize verifies WrTransparentMode_Serialize behavior.
func TestWrTransparentMode_Serialize(t *testing.T) {
	t.Run("serializes write transparent mode command", func(t *testing.T) {
		cmd, _ := NewWrTransparentMode(enums.TransparentModeENABLED)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + TransparentMode(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}
	})
}

// TestNewRdTransparentMode verifies NewRdTransparentMode behavior.
func TestNewRdTransparentMode(t *testing.T) {
	t.Run("creates read transparent mode command", func(t *testing.T) {
		cmd, err := NewRdTransparentMode()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_TRANSPARENT_MODE {
			t.Errorf("expected CommandCode RD_TRANSPARENT_MODE, got 0x%02x", cmd.CommandCode)
		}
	})
}

// TestRdTransparentMode_Serialize verifies RdTransparentMode_Serialize behavior.
func TestRdTransparentMode_Serialize(t *testing.T) {
	t.Run("serializes read transparent mode command", func(t *testing.T) {
		cmd, _ := NewRdTransparentMode()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}
	})
}

// TestParseRdTransparentModeResponseOK verifies ParseRdTransparentModeResponseOK behavior.
func TestParseRdTransparentModeResponseOK(t *testing.T) {
	t.Run("parses transparent mode response", func(t *testing.T) {
		// Response: TransparentMode(1) = 1 byte
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01},
			OptData: nil,
		}

		result, err := ParseRdTransparentModeResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.TransparentMode != enums.TransparentModeENABLED {
			t.Errorf("expected TransparentMode ENABLED, got %v", result.TransparentMode)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01},
			OptData: nil,
		}

		_, err := ParseRdTransparentModeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseRdTransparentModeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}

// TestNewWrTxOnlyMode verifies NewWrTxOnlyMode behavior.
func TestNewWrTxOnlyMode(t *testing.T) {
	t.Run("creates write TX only mode command", func(t *testing.T) {
		cmd, err := NewWrTxOnlyMode(enums.TxOnlyModeENABLED_WITHOUT_AUTO_SLEEP)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandWR_TX_ONLY_MODE {
			t.Errorf("expected CommandCode WR_TX_ONLY_MODE, got 0x%02x", cmd.CommandCode)
		}

		if cmd.TxOnlyMode != enums.TxOnlyModeENABLED_WITHOUT_AUTO_SLEEP {
			t.Errorf("expected TxOnlyMode ENABLED_WITHOUT_AUTO_SLEEP, got %v", cmd.TxOnlyMode)
		}
	})
}

// TestWrTxOnlyMode_Serialize verifies WrTxOnlyMode_Serialize behavior.
func TestWrTxOnlyMode_Serialize(t *testing.T) {
	t.Run("serializes write TX only mode command", func(t *testing.T) {
		cmd, _ := NewWrTxOnlyMode(enums.TxOnlyModeENABLED_WITHOUT_AUTO_SLEEP)
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		// Data: Command(1) + TxOnlyMode(1) = 2 bytes
		if len(telegram.Data) != 2 {
			t.Errorf("expected Data length 2, got %d", len(telegram.Data))
		}
	})
}

// TestNewRdTxOnlyMode verifies NewRdTxOnlyMode behavior.
func TestNewRdTxOnlyMode(t *testing.T) {
	t.Run("creates read TX only mode command", func(t *testing.T) {
		cmd, err := NewRdTxOnlyMode()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cmd.CommandCode != enums.CommonCommandRD_TX_ONLY_MODE {
			t.Errorf("expected CommandCode RD_TX_ONLY_MODE, got 0x%02x", cmd.CommandCode)
		}
	})
}

// TestRdTxOnlyMode_Serialize verifies RdTxOnlyMode_Serialize behavior.
func TestRdTxOnlyMode_Serialize(t *testing.T) {
	t.Run("serializes read TX only mode command", func(t *testing.T) {
		cmd, _ := NewRdTxOnlyMode()
		telegram, err := cmd.Serialize()

		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if len(telegram.Data) != 1 {
			t.Errorf("expected Data length 1, got %d", len(telegram.Data))
		}
	})
}

// TestParseRdTxOnlyModeResponseOK verifies ParseRdTxOnlyModeResponseOK behavior.
func TestParseRdTxOnlyModeResponseOK(t *testing.T) {
	t.Run("parses TX only mode response", func(t *testing.T) {
		// Response: TxOnlyMode(1) = 1 byte
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{0x01},
			OptData: nil,
		}

		result, err := ParseRdTxOnlyModeResponseOK(resp)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if result.TxOnlyMode != enums.TxOnlyModeENABLED_WITHOUT_AUTO_SLEEP {
			t.Errorf("expected TxOnlyMode ENABLED_WITHOUT_AUTO_SLEEP, got %v", result.TxOnlyMode)
		}
	})

	t.Run("returns error for non-success return code", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeERROR,
			Data:    []byte{0x01},
			OptData: nil,
		}

		_, err := ParseRdTxOnlyModeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for non-success return code, got nil")
		}
	})

	t.Run("returns error for insufficient data", func(t *testing.T) {
		resp := response.Packet{
			Code:    enums.ReturnCodeSUCCESS,
			Data:    []byte{},
			OptData: nil,
		}

		_, err := ParseRdTxOnlyModeResponseOK(resp)
		if err == nil {
			t.Fatal("expected error for insufficient data, got nil")
		}

		if err.Error() != "failed to deserialize response" {
			t.Errorf("expected error 'failed to deserialize response', got '%s'", err.Error())
		}
	})
}
