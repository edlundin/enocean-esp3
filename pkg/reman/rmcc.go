package reman

import (
	"encoding/binary"
	"fmt"
)

type ReturnCode byte

const (
	ReturnOK                 ReturnCode = 0x00
	ReturnNotSupported       ReturnCode = 0x01
	ReturnWrongParam         ReturnCode = 0x02
	ReturnOperationDenied    ReturnCode = 0x03
	ReturnSessionClosed      ReturnCode = 0x10
	ReturnInsufficientRights ReturnCode = 0x11
)

func CodePayload(code uint32) ([]byte, error) {
	if code == 0 || code == 0xffffffff {
		return nil, fmt.Errorf("reserved reman code 0x%08x", code)
	}
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, code)
	return b, nil
}

func QueryIDPayload() []byte { return nil }
func PingPayload() []byte    { return nil }

type StatusAnswer struct{ Return ReturnCode }

func ParseStatusAnswer(b []byte) (StatusAnswer, error) {
	if len(b) != 1 {
		return StatusAnswer{}, fmt.Errorf("status answer length %d, want 1", len(b))
	}
	return StatusAnswer{Return: ReturnCode(b[0])}, nil
}
