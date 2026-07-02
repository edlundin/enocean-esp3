package reman

import (
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
)

const (
	MaxPayload      = 508
	MaxParts        = 64
	DefaultStatus   = 0x0f
	ManufacturerID  = 0x07ff
	FuncUnlock      = 0x001
	FuncLock        = 0x002
	FuncSetCode     = 0x003
	FuncQueryID     = 0x004
	FuncPing        = 0x006
	FuncQueryStatus = 0x008
)

type Message struct {
	Seq            byte
	ManufacturerID uint16
	Function       uint16
	Payload        []byte
	SourceID       deviceid.DeviceID
	DestinationID  deviceid.DeviceID
}

func (m Message) Packets() ([]erp1.Packet, error) {
	if m.Seq == 0 || m.Seq > 3 {
		return nil, fmt.Errorf("invalid seq %d", m.Seq)
	}
	if len(m.Payload) > MaxPayload {
		return nil, fmt.Errorf("payload length %d > %d", len(m.Payload), MaxPayload)
	}
	data := append([]byte(nil), m.Payload...)
	first := make([]byte, 8)
	putHeader(first[:4], len(data), m.ManufacturerID, m.Function)
	copy(first[4:], data)
	data = data[min(len(data), 4):]
	parts := [][]byte{first}
	for len(data) > 0 {
		n := min(len(data), 8)
		part := make([]byte, 8)
		copy(part, data[:n])
		parts = append(parts, part)
		data = data[n:]
	}
	if len(parts) > MaxParts {
		return nil, fmt.Errorf("parts %d > %d", len(parts), MaxParts)
	}
	out := make([]erp1.Packet, len(parts))
	for i, part := range parts {
		out[i] = erp1.Packet{Rorg: enums.RorgSYS_EX, UserData: append([]byte{m.Seq<<6 | byte(i)}, part...), SenderID: m.SourceID, DestinationID: m.DestinationID, Status: DefaultStatus, SubTelNum: 1, SecurityLevel: 3, Rssi: 0xff}
	}
	return out, nil
}

type Part struct {
	Seq, Index              byte
	ManufacturerID          uint16
	Function                uint16
	Length                  int
	Payload                 []byte
	SourceID, DestinationID deviceid.DeviceID
}

func ParsePacket(p erp1.Packet) (Part, error) {
	if p.Rorg != enums.RorgSYS_EX {
		return Part{}, fmt.Errorf("not SYS_EX")
	}
	if len(p.UserData) != 9 {
		return Part{}, fmt.Errorf("SYS_EX user data length %d, want 9", len(p.UserData))
	}
	seq, idx := p.UserData[0]>>6, p.UserData[0]&0x3f
	if seq == 0 {
		return Part{}, fmt.Errorf("invalid seq 0")
	}
	part := Part{Seq: seq, Index: idx, SourceID: p.SenderID, DestinationID: p.DestinationID}
	if idx == 0 {
		part.Length, part.ManufacturerID, part.Function = getHeader(p.UserData[1:5])
		part.Payload = append([]byte(nil), p.UserData[5:9]...)
	} else {
		part.Payload = append([]byte(nil), p.UserData[1:9]...)
	}
	return part, nil
}

func Merge(parts []Part) (Message, bool, error) {
	if len(parts) == 0 {
		return Message{}, false, nil
	}
	sort.Slice(parts, func(i, j int) bool { return parts[i].Index < parts[j].Index })
	first := parts[0]
	if first.Index != 0 {
		return Message{}, false, nil
	}
	seen := map[byte]bool{}
	payload := []byte{}
	for _, p := range parts {
		if p.Seq != first.Seq || p.SourceID != first.SourceID || p.DestinationID != first.DestinationID {
			return Message{}, false, fmt.Errorf("mixed message parts")
		}
		if seen[p.Index] {
			return Message{}, false, fmt.Errorf("duplicate index %d", p.Index)
		}
		seen[p.Index] = true
		payload = append(payload, p.Payload...)
	}
	if len(payload) < first.Length {
		return Message{}, false, nil
	}
	return Message{Seq: first.Seq, ManufacturerID: first.ManufacturerID, Function: first.Function, Payload: payload[:first.Length], SourceID: first.SourceID, DestinationID: first.DestinationID}, true, nil
}

func putHeader(b []byte, length int, manufacturer, fn uint16) {
	v := uint32(length&0x1ff)<<23 | uint32(manufacturer&0x7ff)<<12 | uint32(fn&0xfff)
	binary.BigEndian.PutUint32(b, v)
}

func getHeader(b []byte) (int, uint16, uint16) {
	v := binary.BigEndian.Uint32(b)
	return int(v >> 23), uint16((v >> 12) & 0x7ff), uint16(v & 0xfff)
}
