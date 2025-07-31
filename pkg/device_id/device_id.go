package device_id

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

type DeviceID uint32

const (
	DeviceIDSize int = 4
)

func BroadcastId() DeviceID {
	return 0xffffffff
}

func (d DeviceID) ToArray() [DeviceIDSize]byte {
	return [4]byte{
		byte(d >> 24),
		byte(d >> 16),
		byte(d >> 8),
		byte(d),
	}
}

func (d DeviceID) String() string {
	return fmt.Sprintf("%08x", uint32(d))
}

func FromHexString(hexStr string) (DeviceID, error) {
	const sizeMaxStr = DeviceIDSize * 2

	hexStr = strings.TrimPrefix(hexStr, "0x")

	if len(hexStr) > sizeMaxStr {
		return 0, fmt.Errorf("invalid length (got:%d, max:%d)", len(hexStr), sizeMaxStr)
	}

	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	b, err := hex.DecodeString(hexStr)

	if err != nil {
		return 0, err
	}

	return FromByteArray(b)
}

func FromByteArray(b []byte) (DeviceID, error) {
	if len(b) > DeviceIDSize {
		return 0, fmt.Errorf("invalid length (got:%d, need:%d)", len(b), DeviceIDSize)
	}

	if len(b) < DeviceIDSize {
		b = append(make([]byte, DeviceIDSize-len(b)), b...)
	}

	return DeviceID(binary.BigEndian.Uint32(b[0:]) >> (32 - DeviceIDSize*8)), nil
}
