package enocean

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

type DeviceId uint32

func GetBroadcastId() DeviceId {
	return 0xffffffff
}

func (d DeviceId) String() string {
	return fmt.Sprintf("%08x", uint32(d))
}

func DeviceIdFromHexString(s string) (DeviceId, error) {
	b, err := hex.DecodeString(s)

	if err != nil {
		return 0, err
	}

	if len(b) != 4 {
		return 0, errors.New("invalid length")
	}

	return DeviceIdFromByteArray(b), nil
}

func DeviceIdFromByteArray(b []byte) DeviceId {
	return DeviceId(binary.BigEndian.Uint32(b[0:4]))
}
