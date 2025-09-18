package event

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/eep"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

type Event interface {
	Description() enums.EventCode
}

type Packet struct {
	EventCode enums.EventCode
}

func (p Packet) Description() enums.EventCode {
	return p.EventCode
}

type SAReclaimNotSuccessful struct {
	Packet
}

type SAConfirmLearn struct {
	Packet

	PriorityPostmasterCandidate byte
	ManufacturerID              uint16
	EEP                         eep.EEP
	Rssi                        byte
	PostmasterCandidateID       uint32
	SmartACKClientID            uint32
	HopCount                    byte
}

type SALearnAck struct {
	Packet

	ResponseTime uint16
	ConfirmCode  enums.LearnAckConfirmCode
}

type COReady struct {
	Packet

	Cause enums.WakeUpCause
	Mode  enums.WakeUpMode
}

type COEventSecureDevice struct {
	Packet

	Cause    enums.SecureDeviceEventCause
	DeviceID deviceid.DeviceID
}

type CODutyCycleLimit struct {
	Packet

	Cause enums.DutyCycleLimitCause
}

type COTransmitFailed struct {
	Packet

	Cause enums.TransmitFailedCause
}

type COTxDone struct {
	Packet
}

type COLrnModeDisabled struct {
	Packet
}

func NewPacketFromEsp3(telegram esp3.Telegram) (Event, error) {
	if telegram.PacketType != enums.PacketTypeEVENT {
		return Packet{}, errors.New("invalid packet type")
	}

	data := make([]byte, 0, len(telegram.Data)+len(telegram.OptData))
	data = append(data, telegram.Data...)
	data = append(data, telegram.OptData...)

	eventCode, err := enums.ParseEventCodeFromByte(telegram.Data[0])
	if err != nil {
		return Packet{}, err
	}

	switch eventCode {
	case enums.EventCodeSA_RECLAIM_NOT_SUCCESSFUL:
		var p SAReclaimNotSuccessful
		if err := decodeBinary(data, &p); err != nil {
			return Packet{}, err
		}

		return p, nil
	case enums.EventCodeSA_CONFIRM_LEARN:
		var p SAConfirmLearn
		if err := decodeBinary(data, &p); err != nil {
			return Packet{}, err
		}

		if !p.EEP.Rorg.Valid() {
			return Packet{}, errors.New("invalid EEP rorg")
		}

		return p, nil
	case enums.EventCodeSA_LEARN_ACK:
		var p SALearnAck
		if err := decodeBinary(data, &p); err != nil {
			return Packet{}, err
		}

		return p, nil
	case enums.EventCodeCO_READY:
		var p COReady
		if err := decodeBinary(data, &p); err != nil {
			return Packet{}, err
		}

		if !p.Cause.Valid() {
			return Packet{}, errors.New("invalid wake up cause")
		}

		if !p.Mode.Valid() {
			return Packet{}, errors.New("invalid wake up mode")
		}

		return p, nil
	case enums.EventCodeCO_EVENT_SECUREDEVICES:
		var p COEventSecureDevice
		if err := decodeBinary(data, &p); err != nil {
			return Packet{}, err
		}

		if !p.Cause.Valid() {
			return Packet{}, errors.New("invalid secure device event cause")
		}

		return p, nil
	case enums.EventCodeCO_DUTYCYCLE_LIMIT:
		var p CODutyCycleLimit
		if err := decodeBinary(data, &p); err != nil {
			return Packet{}, err
		}

		if !p.Cause.Valid() {
			return Packet{}, errors.New("invalid duty cycle limit cause")
		}

		return p, nil
	case enums.EventCodeCO_TRANSMIT_FAILED:
		var p COTransmitFailed
		if err := decodeBinary(data, &p); err != nil {
			return Packet{}, err
		}

		if !p.Cause.Valid() {
			return Packet{}, errors.New("invalid transmit failed cause")
		}

		return p, nil
	case enums.EventCodeCO_TX_DONE:
		var p COTxDone
		if err := decodeBinary(data, &p); err != nil {
			return Packet{}, err
		}

		return p, nil
	case enums.EventCodeCO_LRN_MODE_DISABLED:
		var p COLrnModeDisabled
		if err := decodeBinary(data, &p); err != nil {
			return Packet{}, err
		}

		return p, nil
	default:
		return Packet{}, errors.New("invalid event code")
	}
}

func decodeBinary(data []byte, out any) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.LittleEndian, out)
}
