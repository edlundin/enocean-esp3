package srm

import (
	"encoding/binary"
	"fmt"
)

const (
	FuncAction            uint16 = 0x005
	FuncPing              uint16 = 0x006
	FuncQueryStatus       uint16 = 0x008
	FuncPingResponse      uint16 = 0x606
	FuncQueryStatusAnswer uint16 = 0x608
	FuncRemoteLearn       uint16 = 0x201
	FuncMemoryWrite       uint16 = 0x203
	FuncMemoryRead        uint16 = 0x204
	FuncMemoryReadAnswer  uint16 = 0x804
	FuncRemoveDevice      uint16 = 0x207
)

type ReturnCode byte

const (
	ReturnOK           ReturnCode = 0x00
	ReturnNotSupported ReturnCode = 0x01
	ReturnWrongParam   ReturnCode = 0x04
	ReturnBusy         ReturnCode = 0x05
	ReturnDenied       ReturnCode = 0x07
	ReturnTimeout      ReturnCode = 0x08
	ReturnNoAck        ReturnCode = 0x0d
	ReturnTooMuchData  ReturnCode = 0x0e
	ReturnWrongState   ReturnCode = 0x0f
)

type Message struct {
	ManufacturerID *uint16
	Function       uint16
	Payload        []byte
}

// MarshalSYSEx marshals SYSEx.
func (m Message) MarshalSYSEx() ([]byte, error) {
	if m.Function == 0 || m.Function > 0x0fff {
		return nil, fmt.Errorf("invalid function 0x%x", m.Function)
	}
	out := make([]byte, 2, 2+len(m.Payload))
	if m.ManufacturerID == nil {
		binary.BigEndian.PutUint16(out[:2], m.Function&0x0fff)
	} else {
		if *m.ManufacturerID > 0x07ff {
			return nil, fmt.Errorf("invalid manufacturer ID 0x%x", *m.ManufacturerID)
		}
		out = make([]byte, 3, 3+len(m.Payload))
		v := uint32(1)<<23 | uint32(*m.ManufacturerID)<<12 | uint32(m.Function&0x0fff)
		out[0], out[1], out[2] = byte(v>>16), byte(v>>8), byte(v)
	}
	out = append(out, m.Payload...)
	return out, nil
}

// ParseSYSEx parses SYSEx.
func ParseSYSEx(b []byte) (Message, error) {
	if len(b) < 2 {
		return Message{}, fmt.Errorf("SYS_EX payload too short")
	}
	if b[0]&0x80 == 0 {
		if b[0]&0x70 != 0 {
			return Message{}, fmt.Errorf("reserved Alliance SYS_EX header bits set")
		}
		fn := binary.BigEndian.Uint16(b[:2]) & 0x0fff
		if fn == 0 {
			return Message{}, fmt.Errorf("invalid function 0x%x", fn)
		}
		return Message{Function: fn, Payload: append([]byte(nil), b[2:]...)}, nil
	}
	if len(b) < 3 {
		return Message{}, fmt.Errorf("manufacturer SYS_EX payload too short")
	}
	v := uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
	mid := uint16((v >> 12) & 0x07ff)
	fn := uint16(v & 0x0fff)
	if fn == 0 {
		return Message{}, fmt.Errorf("invalid function 0x%x", fn)
	}
	return Message{ManufacturerID: &mid, Function: fn, Payload: append([]byte(nil), b[3:]...)}, nil
}

type QueryStatusAnswer struct {
	LastFunction uint16
	Return       ReturnCode
}

// Payload returns the serialized payload.
func (a QueryStatusAnswer) Payload() []byte {
	b := make([]byte, 3)
	binary.BigEndian.PutUint16(b[:2], a.LastFunction&0x0fff)
	b[2] = byte(a.Return)
	return b
}

// ParseQueryStatusAnswer parses QueryStatusAnswer.
func ParseQueryStatusAnswer(b []byte) (QueryStatusAnswer, error) {
	if len(b) != 3 {
		return QueryStatusAnswer{}, fmt.Errorf("query status answer length %d, want 3", len(b))
	}
	if b[0]&0xf0 != 0 {
		return QueryStatusAnswer{}, fmt.Errorf("reserved query status bits set")
	}
	return QueryStatusAnswer{LastFunction: binary.BigEndian.Uint16(b[:2]), Return: ReturnCode(b[2])}, nil
}

type RemoteLearnFlag byte

const (
	RemoteLearnStart            RemoteLearnFlag = 0x01
	RemoteLearnNextChannel      RemoteLearnFlag = 0x02
	RemoteLearnStop             RemoteLearnFlag = 0x03
	RemoteLearnSmartAckSimple   RemoteLearnFlag = 0x04
	RemoteLearnSmartAckAdvanced RemoteLearnFlag = 0x05
	RemoteLearnSmartAckStop     RemoteLearnFlag = 0x06
)

// PingPayload builds a ping payload.
func PingPayload() []byte { return nil }

// PingResponsePayload builds a ping-response payload.
func PingResponsePayload(rssi byte) []byte { return []byte{rssi} }

// RemoteLearnPayload builds a remote-learn payload.
func RemoteLearnPayload(enable bool) []byte {
	if enable {
		return []byte{byte(RemoteLearnStart)}
	}
	return []byte{byte(RemoteLearnStop)}
}

// ParseRemoteLearnPayload parses RemoteLearnPayload.
func ParseRemoteLearnPayload(b []byte) (RemoteLearnFlag, error) {
	if len(b) != 1 || b[0] < byte(RemoteLearnStart) || b[0] > byte(RemoteLearnSmartAckStop) {
		return 0, fmt.Errorf("invalid remote learn payload")
	}
	return RemoteLearnFlag(b[0]), nil
}

// MemoryReadPayload builds a memory-read payload.
func MemoryReadPayload(addr uint32, n byte) []byte {
	b := make([]byte, 5)
	binary.BigEndian.PutUint32(b[:4], addr)
	b[4] = n
	return b
}

// MemoryWritePayload builds a memory-write payload.
func MemoryWritePayload(addr uint32, data []byte) ([]byte, error) {
	if len(data) > 255 {
		return nil, fmt.Errorf("memory write length %d > 255", len(data))
	}
	b := make([]byte, 5, 5+len(data))
	binary.BigEndian.PutUint32(b[:4], addr)
	b[4] = byte(len(data))
	return append(b, data...), nil
}
