package recom

import (
	"encoding/binary"
	"fmt"

	"github.com/edlundin/enocean-esp3/pkg/deviceid"
)

const ManufacturerID uint16 = 0x07ff

const (
	FuncGetLinkTableMetadata uint16 = 0x210
	FuncGetLinkTable         uint16 = 0x211
	FuncSetLinkTableContent  uint16 = 0x212
	FuncGetProductID         uint16 = 0x227
	FuncProductIDResponse    uint16 = 0x827
	FuncAcknowledge          uint16 = 0x240
)

type Direction byte

const (
	Inbound  Direction = 0
	Outbound Direction = 1
)

type ProductID struct {
	Manufacturer uint16
	Product      uint32
}

func (p ProductID) MarshalBinary() []byte {
	b := make([]byte, 6)
	binary.BigEndian.PutUint16(b[:2], p.Manufacturer)
	binary.BigEndian.PutUint32(b[2:], p.Product)
	return b
}

func ParseProductID(b []byte) (ProductID, error) {
	if len(b) != 6 {
		return ProductID{}, fmt.Errorf("product ID length %d, want 6", len(b))
	}
	return ProductID{Manufacturer: binary.BigEndian.Uint16(b[:2]), Product: binary.BigEndian.Uint32(b[2:])}, nil
}

type ParamRecord struct {
	Index uint16
	Value []byte
}

func MarshalParamRecords(records []ParamRecord) ([]byte, error) {
	var out []byte
	for _, r := range records {
		if len(r.Value) > 64 {
			return nil, fmt.Errorf("parameter 0x%04x length %d > 64", r.Index, len(r.Value))
		}
		out = binary.BigEndian.AppendUint16(out, r.Index)
		out = append(out, byte(len(r.Value)))
		out = append(out, r.Value...)
	}
	if len(out) > 67 {
		return nil, fmt.Errorf("parameter payload length %d > 67", len(out))
	}
	return out, nil
}

func ParseParamRecords(b []byte) ([]ParamRecord, error) {
	var records []ParamRecord
	for len(b) > 0 {
		if len(b) < 3 {
			return nil, fmt.Errorf("truncated parameter record")
		}
		idx, n := binary.BigEndian.Uint16(b[:2]), int(b[2])
		b = b[3:]
		if n > 64 || len(b) < n {
			return nil, fmt.Errorf("invalid parameter 0x%04x length %d", idx, n)
		}
		records = append(records, ParamRecord{Index: idx, Value: append([]byte(nil), b[:n]...)})
		b = b[n:]
	}
	return records, nil
}

type LinkEntry struct {
	EEP      [3]byte
	DeviceID deviceid.DeviceID
	Data     [2]byte
}

func (e LinkEntry) MarshalBinary() []byte {
	b := make([]byte, 9)
	copy(b[:3], e.EEP[:])
	id := e.DeviceID.ToArray()
	copy(b[3:7], id[:])
	copy(b[7:], e.Data[:])
	return b
}

func ParseLinkEntry(b []byte) (LinkEntry, error) {
	if len(b) != 9 {
		return LinkEntry{}, fmt.Errorf("link entry length %d, want 9", len(b))
	}
	id, _ := deviceid.FromByteArray(b[3:7])
	return LinkEntry{EEP: [3]byte{b[0], b[1], b[2]}, DeviceID: id, Data: [2]byte{b[7], b[8]}}, nil
}
