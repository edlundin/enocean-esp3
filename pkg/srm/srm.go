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

func (m Message) MarshalSYSEx() ([]byte, error) {
	if m.Function == 0 || m.Function > 0x0fff {
		return nil, fmt.Errorf("invalid function 0x%x", m.Function)
	}
	out := make([]byte, 2, 2+len(m.Payload))
	if m.ManufacturerID == nil {
		binary.BigEndian.PutUint16(out[:2], m.Function&0x0fff)
	} else {
		out = make([]byte, 3, 3+len(m.Payload))
		v := uint32(1)<<23 | uint32(*m.ManufacturerID&0x07ff)<<12 | uint32(m.Function&0x0fff)
		out[0], out[1], out[2] = byte(v>>16), byte(v>>8), byte(v)
	}
	out = append(out, m.Payload...)
	return out, nil
}

func ParseSYSEx(b []byte) (Message, error) {
	if len(b) < 2 {
		return Message{}, fmt.Errorf("SYS_EX payload too short")
	}
	if b[0]&0x80 == 0 {
		fn := binary.BigEndian.Uint16(b[:2]) & 0x0fff
		return Message{Function: fn, Payload: append([]byte(nil), b[2:]...)}, nil
	}
	if len(b) < 3 {
		return Message{}, fmt.Errorf("manufacturer SYS_EX payload too short")
	}
	v := uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
	mid := uint16((v >> 12) & 0x07ff)
	return Message{ManufacturerID: &mid, Function: uint16(v & 0x0fff), Payload: append([]byte(nil), b[3:]...)}, nil
}

type QueryStatusAnswer struct {
	LastFunction uint16
	Return       ReturnCode
}

func (a QueryStatusAnswer) Payload() []byte {
	b := make([]byte, 3)
	binary.BigEndian.PutUint16(b[:2], a.LastFunction&0x0fff)
	b[2] = byte(a.Return)
	return b
}

func ParseQueryStatusAnswer(b []byte) (QueryStatusAnswer, error) {
	if len(b) != 3 {
		return QueryStatusAnswer{}, fmt.Errorf("query status answer length %d, want 3", len(b))
	}
	return QueryStatusAnswer{LastFunction: binary.BigEndian.Uint16(b[:2]) & 0x0fff, Return: ReturnCode(b[2])}, nil
}

func PingPayload() []byte                  { return nil }
func PingResponsePayload(rssi byte) []byte { return []byte{rssi} }
func RemoteLearnPayload(enable bool) []byte {
	if enable {
		return []byte{1}
	}
	return []byte{0}
}

func MemoryReadPayload(addr uint32, n byte) []byte {
	b := make([]byte, 5)
	binary.BigEndian.PutUint32(b[:4], addr)
	b[4] = n
	return b
}

func MemoryWritePayload(addr uint32, data []byte) ([]byte, error) {
	if len(data) > 255 {
		return nil, fmt.Errorf("memory write length %d > 255", len(data))
	}
	b := make([]byte, 5, 5+len(data))
	binary.BigEndian.PutUint32(b[:4], addr)
	b[4] = byte(len(data))
	return append(b, data...), nil
}
