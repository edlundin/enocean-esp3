package event

import (
	"encoding/binary"
	"testing"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

// TestNewPacketFromEsp3 verifies NewPacketFromEsp3 behavior.
func TestNewPacketFromEsp3(t *testing.T) {
	t.Run("returns error for invalid packet type", func(t *testing.T) {
		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeRADIO_ERP1,
			Data:       []byte{0x01},
			OptData:    []byte{},
		}

		_, err := NewPacketFromEsp3(telegram)
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		expectedError := "invalid packet type"
		if err.Error() != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
		}
	})

	t.Run("returns error for invalid event code", func(t *testing.T) {
		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       []byte{0xFF}, // Invalid event code
			OptData:    []byte{},
		}

		_, err := NewPacketFromEsp3(telegram)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("parses SA_RECLAIM_NOT_SUCCESSFUL event", func(t *testing.T) {
		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       []byte{byte(enums.EventCodeSA_RECLAIM_NOT_SUCCESSFUL)},
			OptData:    []byte{},
		}

		event, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if event.Description() != enums.EventCodeSA_RECLAIM_NOT_SUCCESSFUL {
			t.Errorf("expected event code SA_RECLAIM_NOT_SUCCESSFUL, got %v", event.Description())
		}

		_, ok := event.(SAReclaimNotSuccessful)
		if !ok {
			t.Errorf("expected SAReclaimNotSuccessful type, got %T", event)
		}
	})

	t.Run("parses SA_CONFIRM_LEARN event", func(t *testing.T) {
		// SA_CONFIRM_LEARN structure:
		// EventCode (1) + PriorityPostmasterCandidate (1) + ManufacturerID (2) + EEP (3) + RSSI (1) + PostmasterCandidateID (4) + SmartACKClientID (4) + HopCount (1) = 17 bytes
		data := make([]byte, 17)
		data[0] = byte(enums.EventCodeSA_CONFIRM_LEARN)
		data[1] = 0x01                                         // PriorityPostmasterCandidate
		binary.BigEndian.PutUint16(data[2:4], 0x1234)       // ManufacturerID
		data[4] = byte(enums.RorgVLD)                          // EEP Rorg
		data[5] = 0x02                                         // EEP Func
		data[6] = 0x03                                         // EEP Type
		data[7] = 0x80                                         // RSSI
		binary.BigEndian.PutUint32(data[8:12], 0x12345678)  // PostmasterCandidateID
		binary.BigEndian.PutUint32(data[12:16], 0x87654321) // SmartACKClientID
		data[16] = 0x05                                        // HopCount

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		event, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if event.Description() != enums.EventCodeSA_CONFIRM_LEARN {
			t.Errorf("expected event code SA_CONFIRM_LEARN, got %v", event.Description())
		}

		saConfirmLearn, ok := event.(SAConfirmLearn)
		if !ok {
			t.Errorf("expected SAConfirmLearn type, got %T", event)
		}

		if saConfirmLearn.PriorityPostmasterCandidate != 0x01 {
			t.Errorf("expected PriorityPostmasterCandidate 0x01, got 0x%02x", saConfirmLearn.PriorityPostmasterCandidate)
		}
		if saConfirmLearn.ManufacturerID != 0x1234 {
			t.Errorf("expected ManufacturerID 0x1234, got 0x%04x", saConfirmLearn.ManufacturerID)
		}
		if saConfirmLearn.EEP.Rorg != enums.RorgVLD {
			t.Errorf("expected EEP Rorg VLD, got %v", saConfirmLearn.EEP.Rorg)
		}
		if saConfirmLearn.Rssi != 0x80 {
			t.Errorf("expected RSSI 0x80, got 0x%02x", saConfirmLearn.Rssi)
		}
		if saConfirmLearn.PostmasterCandidateID != 0x12345678 {
			t.Errorf("expected PostmasterCandidateID 0x12345678, got 0x%08x", saConfirmLearn.PostmasterCandidateID)
		}
		if saConfirmLearn.SmartACKClientID != 0x87654321 {
			t.Errorf("expected SmartACKClientID 0x87654321, got 0x%08x", saConfirmLearn.SmartACKClientID)
		}
		if saConfirmLearn.HopCount != 0x05 {
			t.Errorf("expected HopCount 0x05, got 0x%02x", saConfirmLearn.HopCount)
		}
	})

	t.Run("returns error for SA_CONFIRM_LEARN with invalid EEP", func(t *testing.T) {
		data := make([]byte, 17)
		data[0] = byte(enums.EventCodeSA_CONFIRM_LEARN)
		data[4] = 0xFF // Invalid Rorg

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		_, err := NewPacketFromEsp3(telegram)
		if err == nil {
			t.Errorf("expected error for invalid EEP, got nil")
		}
	})

	t.Run("parses SA_LEARN_ACK event", func(t *testing.T) {
		// SA_LEARN_ACK structure: EventCode (1) + ResponseTime (2) + ConfirmCode (1) = 4 bytes
		data := make([]byte, 4)
		data[0] = byte(enums.EventCodeSA_LEARN_ACK)
		binary.BigEndian.PutUint16(data[1:3], 0x1234) // ResponseTime
		data[3] = byte(enums.LearnAckConfirmCodeLRN_IN)  // ConfirmCode

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		event, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		saLearnAck, ok := event.(SALearnAck)
		if !ok {
			t.Errorf("expected SALearnAck type, got %T", event)
		}

		if saLearnAck.ResponseTime != 0x1234 {
			t.Errorf("expected ResponseTime 0x1234, got 0x%04x", saLearnAck.ResponseTime)
		}
		if saLearnAck.ConfirmCode != enums.LearnAckConfirmCodeLRN_IN {
			t.Errorf("expected ConfirmCode LRN_IN, got %v", saLearnAck.ConfirmCode)
		}
	})

	t.Run("parses CO_READY event", func(t *testing.T) {
		// CO_READY structure: EventCode (1) + Cause (1) + Mode (1) = 3 bytes
		data := make([]byte, 3)
		data[0] = byte(enums.EventCodeCO_READY)
		data[1] = byte(enums.WakeUpCauseUART_WAKE_UP)     // Cause
		data[2] = byte(enums.WakeUpModeEXTENDED_SECURITY) // Mode

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		event, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		coReady, ok := event.(COReady)
		if !ok {
			t.Errorf("expected COReady type, got %T", event)
		}

		if coReady.Cause != enums.WakeUpCauseUART_WAKE_UP {
			t.Errorf("expected Cause UART_WAKE_UP, got %v", coReady.Cause)
		}
		if coReady.Mode != enums.WakeUpModeEXTENDED_SECURITY {
			t.Errorf("expected Mode EXTENDED_SECURITY, got %v", coReady.Mode)
		}
	})

	t.Run("returns error for CO_READY with invalid cause", func(t *testing.T) {
		data := make([]byte, 3)
		data[0] = byte(enums.EventCodeCO_READY)
		data[1] = 0xFF // Invalid cause
		data[2] = byte(enums.WakeUpModeSTANDARD_SECURITY)

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		_, err := NewPacketFromEsp3(telegram)
		if err == nil {
			t.Errorf("expected error for invalid cause, got nil")
		}
	})

	t.Run("returns error for CO_READY with invalid mode", func(t *testing.T) {
		data := make([]byte, 3)
		data[0] = byte(enums.EventCodeCO_READY)
		data[1] = byte(enums.WakeUpCauseUART_WAKE_UP)
		data[2] = 0xFF // Invalid mode

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		_, err := NewPacketFromEsp3(telegram)
		if err == nil {
			t.Errorf("expected error for invalid mode, got nil")
		}
	})

	t.Run("parses CO_EVENT_SECUREDEVICES event", func(t *testing.T) {
		// CO_EVENT_SECUREDEVICES structure: EventCode (1) + Cause (1) + DeviceID (4) = 6 bytes
		// Note: decodeBinary uses LittleEndian, so DeviceID must be written in little-endian order
		data := make([]byte, 6)
		data[0] = byte(enums.EventCodeCO_EVENT_SECUREDEVICES)
		data[1] = byte(enums.COEventSecureTEACH_IN_SUCCESSFUL) // Cause
		deviceID := deviceid.DeviceID(0x12345678)
		// Write DeviceID in little-endian order for binary.BigEndian.Read
		binary.BigEndian.PutUint32(data[2:6], uint32(deviceID))

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		event, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		coEventSecure, ok := event.(COEventSecureDevice)
		if !ok {
			t.Errorf("expected COEventSecureDevice type, got %T", event)
		}

		if coEventSecure.Cause != enums.COEventSecureTEACH_IN_SUCCESSFUL {
			t.Errorf("expected Cause TEACH_IN_SUCCESSFUL, got %v", coEventSecure.Cause)
		}
		if coEventSecure.DeviceID != deviceID {
			t.Errorf("expected DeviceID %v, got %v", deviceID, coEventSecure.DeviceID)
		}
	})

	t.Run("parses CO_DUTYCYCLE_LIMIT event", func(t *testing.T) {
		// CO_DUTYCYCLE_LIMIT structure: EventCode (1) + Cause (1) = 2 bytes
		data := make([]byte, 2)
		data[0] = byte(enums.EventCodeCO_DUTYCYCLE_LIMIT)
		data[1] = byte(enums.DutyCycleLimitCauseREACHED) // Cause

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		event, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		coDutyCycleLimit, ok := event.(CODutyCycleLimit)
		if !ok {
			t.Errorf("expected CODutyCycleLimit type, got %T", event)
		}

		if coDutyCycleLimit.Cause != enums.DutyCycleLimitCauseREACHED {
			t.Errorf("expected Cause REACHED, got %v", coDutyCycleLimit.Cause)
		}
	})

	t.Run("parses CO_TRANSMIT_FAILED event", func(t *testing.T) {
		// CO_TRANSMIT_FAILED structure: EventCode (1) + Cause (1) = 2 bytes
		data := make([]byte, 2)
		data[0] = byte(enums.EventCodeCO_TRANSMIT_FAILED)
		data[1] = byte(enums.TransmitFailedCauseCSMA_FAILED_CHANNEL_NOT_FREE) // Cause

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		event, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		coTransmitFailed, ok := event.(COTransmitFailed)
		if !ok {
			t.Errorf("expected COTransmitFailed type, got %T", event)
		}

		if coTransmitFailed.Cause != enums.TransmitFailedCauseCSMA_FAILED_CHANNEL_NOT_FREE {
			t.Errorf("expected Cause CSMA_FAILED_CHANNEL_NOT_FREE, got %v", coTransmitFailed.Cause)
		}
	})

	t.Run("parses CO_TX_DONE event", func(t *testing.T) {
		// CO_TX_DONE structure: EventCode (1) = 1 byte
		data := make([]byte, 1)
		data[0] = byte(enums.EventCodeCO_TX_DONE)

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		event, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		_, ok := event.(COTxDone)
		if !ok {
			t.Errorf("expected COTxDone type, got %T", event)
		}
	})

	t.Run("parses CO_LRN_MODE_DISABLED event", func(t *testing.T) {
		// CO_LRN_MODE_DISABLED structure: EventCode (1) = 1 byte
		data := make([]byte, 1)
		data[0] = byte(enums.EventCodeCO_LRN_MODE_DISABLED)

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    []byte{},
		}

		event, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		_, ok := event.(COLrnModeDisabled)
		if !ok {
			t.Errorf("expected COLrnModeDisabled type, got %T", event)
		}
	})

	t.Run("handles optdata in event parsing", func(t *testing.T) {
		// Events can have optdata that gets concatenated with data
		data := []byte{byte(enums.EventCodeCO_TX_DONE)}
		optData := []byte{0x01, 0x02, 0x03}

		telegram := esp3.Telegram{
			PacketType: enums.PacketTypeEVENT,
			Data:       data,
			OptData:    optData,
		}

		event, err := NewPacketFromEsp3(telegram)
		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if event.Description() != enums.EventCodeCO_TX_DONE {
			t.Errorf("expected event code CO_TX_DONE, got %v", event.Description())
		}
	})

	t.Run("handles all learn ack confirm codes", func(t *testing.T) {
		confirmCodes := []struct {
			code enums.LearnAckConfirmCode
			name string
		}{
			{enums.LearnAckConfirmCodeLRN_IN, "LRN_IN"},
			{enums.LearnAckConfirmCodeEEP_NOT_ACCEPTED, "EEP_NOT_ACCEPTED"},
			{enums.LearnAckConfirmCodeNO_PLACE_IN_PM, "NO_PLACE_IN_PM"},
			{enums.LearnAckConfirmCodeNO_PLACE_IN_CONTROLLER, "NO_PLACE_IN_CONTROLLER"},
			{enums.LearnAckConfirmCodeRSSI_NOT_GOOD_ENOUGH, "RSSI_NOT_GOOD_ENOUGH"},
			{enums.LearnAckConfirmCodeLRN_OUT, "LRN_OUT"},
			{enums.LearnAckConfirmCodeFUNCTION_NOT_SUPPORTED, "FUNCTION_NOT_SUPPORTED"},
		}

		for _, tc := range confirmCodes {
			t.Run(tc.name, func(t *testing.T) {
				data := make([]byte, 4)
				data[0] = byte(enums.EventCodeSA_LEARN_ACK)
				binary.BigEndian.PutUint16(data[1:3], 0x0000)
				data[3] = byte(tc.code)

				telegram := esp3.Telegram{
					PacketType: enums.PacketTypeEVENT,
					Data:       data,
					OptData:    []byte{},
				}

				event, err := NewPacketFromEsp3(telegram)
				if err != nil {
					t.Errorf("failed to parse %s: %v", tc.name, err)
					return
				}

				saLearnAck, ok := event.(SALearnAck)
				if !ok {
					t.Errorf("expected SALearnAck type, got %T", event)
					return
				}

				if saLearnAck.ConfirmCode != tc.code {
					t.Errorf("expected ConfirmCode %v, got %v", tc.code, saLearnAck.ConfirmCode)
				}
			})
		}
	})
}
