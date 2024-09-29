package pkg

import (
	"fmt"
	"testing"

	"github.com/shoenig/test"
)

func TestDeviceId_GetBroadcastId(t *testing.T) {
	t.Run("returns 0xffffffff (broadcast ID)", func(t *testing.T) {
		test.Eq(t, DeviceId(0xffffffff), GetBroadcastId())
	})
}

func TestDeviceId_String(t *testing.T) {
	t.Run("returns lower case without hex prefix", func(t *testing.T) {
		deviceId := DeviceId(0xffabc123)
		test.Eq(t, "ffabc123", deviceId.String())
	})

	t.Run("always returns 8 characters, padded with 0", func(t *testing.T) {
		deviceId := DeviceId(0)
		test.Eq(t, "00000000", deviceId.String())

		deviceId = DeviceId(0x42)
		test.Eq(t, "00000042", deviceId.String())

		deviceId = DeviceId(0xabcdef42)
		test.Eq(t, "abcdef42", deviceId.String())
	})
}

func TestDeviceIdFromHexString(t *testing.T) {
	const SIZE_MAX_STR = SIZE_DEVICE_ID * 2

	t.Run("returns the DeviceId representation of the string content (even length)", func(t *testing.T) {
		deviceId, err := DeviceIdFromHexString("ffabcdef")

		test.NoError(t, err)
		test.Eq(t, DeviceId(0xffabcdef), deviceId)
	})

	t.Run("return the DeviceId representation of the string content (odd length with padding)", func(t *testing.T) {
		deviceId, err := DeviceIdFromHexString("0xffa")

		test.NoError(t, err)
		test.Eq(t, DeviceId(0xffa), deviceId)
	})

	t.Run("ignore 0x prefix if present", func(t *testing.T) {
		idStr := "0xffabcdef"
		deviceId, err := DeviceIdFromHexString(idStr)

		test.NoError(t, err)
		test.Eq(t, DeviceId(0xffabcdef), deviceId)
	})

	t.Run(fmt.Sprintf("return an error when the string's length is greater than %d", SIZE_MAX_STR), func(t *testing.T) {
		idStr := "ffabcdefaa"
		_, err := DeviceIdFromHexString(idStr)

		test.EqError(t, err, fmt.Sprintf("invalid length (got:%d, max:%d)", len(idStr), SIZE_MAX_STR))
	})

	t.Run("return an error when the string is invalid (not in hexadecimal format)", func(t *testing.T) {
		idStr := "ffabcdeg"
		_, err := DeviceIdFromHexString(idStr)

		test.Error(t, err)
	})
}

func TestDeviceIdFromByteArray(t *testing.T) {
	t.Run("returns the DeviceId representation of the array content", func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd, 0xef}
		deviceId, err := DeviceIdFromByteArray(byteArray)

		test.NoError(t, err)
		test.Eq(t, DeviceId(0xffabcdef), deviceId)
	})

	t.Run(fmt.Sprintf("return an error when the array's length is greater than %d", SIZE_DEVICE_ID), func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd, 0xef, 0xaa}
		_, err := DeviceIdFromByteArray(byteArray)

		test.EqError(t, err, fmt.Sprintf("invalid length (got:%d, need:%d)", len(byteArray), SIZE_DEVICE_ID))
	})

	t.Run(fmt.Sprintf("return the padded DeviceId when the array's length is lesser than %d", SIZE_DEVICE_ID), func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd}
		deviceId, err := DeviceIdFromByteArray(byteArray)

		test.NoError(t, err)
		test.Eq(t, DeviceId(0xffabcd), deviceId)
	})
}
