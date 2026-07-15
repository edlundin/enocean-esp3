package security

import (
	"errors"
	"fmt"
	"sort"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/erp1"
)

const (
	MaxChainParts = 64
	MaxChainData  = 10 + 13*(MaxChainParts-1)
)

type ChainPart struct {
	Seq, Index byte
	Length     int
	Data       []byte
}

// SplitSEC_CDM splits SEC_CDM.
func SplitSEC_CDM(seq byte, data []byte) ([]erp1.Packet, error) {
	if seq == 0 || seq > 3 {
		return nil, fmt.Errorf("invalid seq %d", seq)
	}
	if len(data) > MaxChainData {
		return nil, fmt.Errorf("chain data too long")
	}
	first := make([]byte, 2, 18)
	first[0], first[1] = byte(len(data)>>8), byte(len(data))
	n := min(len(data), 10)
	first = append(first, data[:n]...)
	data = data[n:]
	parts := [][]byte{first}
	for len(data) > 0 {
		n = min(len(data), 13)
		parts = append(parts, append([]byte(nil), data[:n]...))
		data = data[n:]
	}
	out := make([]erp1.Packet, len(parts))
	for i, p := range parts {
		out[i] = erp1.Packet{Rorg: enums.RorgSEC_CDM, UserData: append([]byte{seq<<6 | byte(i)}, p...), SecurityLevel: RLCExplicit32CMAC32VAES}
	}
	return out, nil
}

// ParseSEC_CDM parses SEC_CDM.
func ParseSEC_CDM(p erp1.Packet) (ChainPart, error) {
	if p.Rorg != enums.RorgSEC_CDM {
		return ChainPart{}, errors.New("not SEC_CDM")
	}
	if len(p.UserData) < 2 || len(p.UserData) > 14 {
		return ChainPart{}, errors.New("invalid SEC_CDM length")
	}
	seq, idx := p.UserData[0]>>6, p.UserData[0]&0x3f
	if seq == 0 {
		return ChainPart{}, errors.New("invalid seq 0")
	}
	part := ChainPart{Seq: seq, Index: idx}
	if idx == 0 {
		if len(p.UserData) < 3 {
			return ChainPart{}, errors.New("SEC_CDM first part too short")
		}
		part.Length = int(p.UserData[1])<<8 | int(p.UserData[2])
		if part.Length > MaxChainData || len(p.UserData) > 13 {
			return ChainPart{}, errors.New("invalid SEC_CDM first part")
		}
		part.Data = append([]byte(nil), p.UserData[3:]...)
	} else {
		part.Data = append([]byte(nil), p.UserData[1:]...)
	}
	return part, nil
}

// MergeSEC_CDM merges SEC_CDM.
func MergeSEC_CDM(parts []ChainPart) ([]byte, bool, error) {
	if len(parts) == 0 {
		return nil, false, nil
	}
	sort.Slice(parts, func(i, j int) bool { return parts[i].Index < parts[j].Index })
	first := parts[0]
	if first.Index != 0 {
		return nil, false, nil
	}
	seen := map[byte]bool{}
	var data []byte
	for i, p := range parts {
		if p.Seq != first.Seq {
			return nil, false, errors.New("mixed seq")
		}
		if seen[p.Index] {
			return nil, false, fmt.Errorf("duplicate index %d", p.Index)
		}
		seen[p.Index] = true
		if p.Index != byte(i) {
			return nil, false, nil
		}
		data = append(data, p.Data...)
	}
	if len(data) < first.Length {
		return nil, false, nil
	}
	if len(data) > first.Length {
		return nil, false, errors.New("SEC_CDM data exceeds declared length")
	}
	return data, true, nil
}
