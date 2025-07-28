package pkg

import device_id "github.com/edlundin/enocean-esp3/pkg/device-id"

type EventCode uint8

const (
	EVENT_CODE_SA_RECLAIM_NOT_SUCCESSFUL EventCode = 0x01
	EVENT_CODE_SA_CONFIRM_LEARN          EventCode = 0x02
	EVENT_CODE_SA_LEARN_ACK              EventCode = 0x03
	EVENT_CODE_CO_READY                  EventCode = 0x04
	EVENT_CODE_CO_EVENT_SECUREDEVICES    EventCode = 0x05
	EVENT_CODE_CO_DUTYCYCLE_LIMIT        EventCode = 0x06
	EVENT_CODE_CO_TRANSMIT_FAILED        EventCode = 0x07
	EVENT_CODE_CO_TX_DONE                EventCode = 0x08
	EVENT_CODE_CO_LRN_MODE_DISABLED      EventCode = 0x09
)

func (eventCode EventCode) String() string {
	switch eventCode {
	case EVENT_CODE_SA_RECLAIM_NOT_SUCCESSFUL:
		return "SA_RECLAIM_NOT_SUCCESSFUL"
	case EVENT_CODE_SA_CONFIRM_LEARN:
		return "SA_CONFIRM_LEARN"
	case EVENT_CODE_SA_LEARN_ACK:
		return "SA_LEARN_ACK"
	case EVENT_CODE_CO_READY:
		return "CO_READY"
	case EVENT_CODE_CO_EVENT_SECUREDEVICES:
		return "CO_EVENT_SECUREDEVICES"
	case EVENT_CODE_CO_DUTYCYCLE_LIMIT:
		return "CO_DUTYCYCLE_LIMIT"
	case EVENT_CODE_CO_TRANSMIT_FAILED:
		return "CO_TRANSMIT_FAILED"
	case EVENT_CODE_CO_TX_DONE:
		return "CO_TX_DONE"
	case EVENT_CODE_CO_LRN_MODE_DISABLED:
		return "CO_LRN_MODE_DISABLED"
	default:
		return "UNKNOWN"
	}
}

//type EventParser interface {
//	Parse(esp3Telegram Esp3Telegram) (Event, error)
//}

type EventPacket struct {
	EventCode EventCode
}

type EventSAReclaimNotSuccessful EventPacket

type EventSAConfirmLearn struct {
	EventPacket

	/* PriorityPostMasterCandidate
	 * 0x08: already post master
	 * 0x04: place in mailbox
	 * 0x02: good RSSI
	 * 0x01: Local
	 */
	PriorityPostmasterCandidate uint8
	ManufacturerID              uint8
	EEP                         device_id.DeviceID
	RSSI                        uint8
	PostmasterCandidateID       uint32
	SmartACKClientID            uint32
	HopCount                    uint8
}
