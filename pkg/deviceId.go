package pkg

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

type DeviceID uint32

const (
	sizeDeviceID int = 4
)

func GetBroadcastId() DeviceID {
	return 0xffffffff
}

func (d DeviceID) String() string {
	return fmt.Sprintf("%08x", uint32(d))
}

func DeviceIdFromHexString(s string) (DeviceID, error) {
	const sizeMaxStr = sizeDeviceID * 2

	stringToDecode := s

	if strings.HasPrefix(stringToDecode, "0x") {
		stringToDecode = s[2:]
	}

	if len(stringToDecode) > sizeMaxStr {
		return 0, fmt.Errorf("invalid length (got:%d, max:%d)", len(stringToDecode), sizeMaxStr)
	}

	if len(stringToDecode)%2 != 0 {
		stringToDecode = strings.Repeat("0", 1) + stringToDecode
	}

	b, err := hex.DecodeString(stringToDecode)

	if err != nil {
		return 0, err
	}

	return DeviceIdFromByteArray(b)
}

func DeviceIdFromByteArray(b []byte) (DeviceID, error) {
	if len(b) > sizeDeviceID {
		return 0, fmt.Errorf("invalid length (got:%d, need:%d)", len(b), sizeDeviceID)
	}

	if len(b) < sizeDeviceID {
		b = append(make([]byte, sizeDeviceID-len(b)), b...)
	}

	return DeviceID(binary.BigEndian.Uint32(b[0:]) >> (32 - sizeDeviceID*8)), nil
}
