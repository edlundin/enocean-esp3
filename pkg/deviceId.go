package pkg

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

type DeviceId uint32

const (
	SIZE_DEVICE_ID int = 4
)

func GetBroadcastId() DeviceId {
	return 0xffffffff
}

func (d DeviceId) String() string {
	return fmt.Sprintf("%08x", uint32(d))
}

func DeviceIdFromHexString(s string) (DeviceId, error) {
	const SIZE_MAX_STR = SIZE_DEVICE_ID * 2

	stringToDecode := s

	if strings.HasPrefix(stringToDecode, "0x") {
		stringToDecode = s[2:]
	}

	if len(stringToDecode) > SIZE_MAX_STR {
		return 0, fmt.Errorf("invalid length (got:%d, max:%d)", len(stringToDecode), SIZE_MAX_STR)
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

func DeviceIdFromByteArray(b []byte) (DeviceId, error) {
	if len(b) > SIZE_DEVICE_ID {
		return 0, fmt.Errorf("invalid length (got:%d, need:%d)", len(b), SIZE_DEVICE_ID)
	}

	if len(b) < SIZE_DEVICE_ID {
		b = append(make([]byte, SIZE_DEVICE_ID-len(b)), b...)
	}

	return DeviceId(binary.BigEndian.Uint32(b[0:]) >> (32 - SIZE_DEVICE_ID*8)), nil
}
