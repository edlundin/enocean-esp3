package enocean

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceId_GetBroadcastId(t *testing.T) {
	t.Run("returns 0xffffffff (broadcast ID)", func(t *testing.T) {
		assert.Equal(t, DeviceId(0xffffffff), GetBroadcastId())
	})
}

func TestDeviceId_String(t *testing.T) {
	t.Run("returns lower case without hex prefix", func(t *testing.T) {
		deviceId := DeviceId(0xffabc123)
		assert.Equal(t, "ffabc123", deviceId.String())
	})

	t.Run("always returns 8 characters, padded with 0", func(t *testing.T) {
		deviceId := DeviceId(0)
		assert.Equal(t, "00000000", deviceId.String())

		deviceId = DeviceId(0x42)
		assert.Equal(t, "00000042", deviceId.String())

		deviceId = DeviceId(0xabcdef42)
		assert.Equal(t, "abcdef42", deviceId.String())
	})
}

func TestDeviceId_DeviceIdFromHexString(t *testing.T) {
	t.Run("returns the DeviceId representation of the string content (even length)", func(t *testing.T) {
		deviceId, err := DeviceIdFromHexString("ffabcdef")

		assert.NoError(t, err)
		assert.Equal(t, DeviceId(0xffabcdef), deviceId)
	})

	t.Run("return the DeviceId representation of the string content (odd length with padding)", func(t *testing.T) {
		deviceId, err := DeviceIdFromHexString("0xffa")

		assert.NoError(t, err)
		assert.Equal(t, DeviceId(0xffa), deviceId)
	})

	t.Run("ignore 0x prefix if present", func(t *testing.T) {
		idStr := "0xffabcdef"
		deviceId, err := DeviceIdFromHexString(idStr)

		assert.NoError(t, err)
		assert.Equal(t, DeviceId(0xffabcdef), deviceId)
	})

	t.Run(fmt.Sprintf("return an error when the string's length is grater than %d", SIZE_DEVICE_ID*2), func(t *testing.T) {
		idStr := "ffabcdefaa"
		_, err := DeviceIdFromHexString(idStr)

		assert.Errorf(t, err, "invalid length (got %d, max: %d)", len(idStr), SIZE_DEVICE_ID)
	})

	t.Run("return an error when the string is invalid (not in hexadecimal format)", func(t *testing.T) {
		idStr := "ffabcdeg"
		_, err := DeviceIdFromHexString(idStr)

		assert.Error(t, err)
	})
}

func TestDeviceId_DeviceIdFromByteArray(t *testing.T) {
	t.Run("returns the DeviceId representation of the array content", func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd, 0xef}
		deviceId, err := DeviceIdFromByteArray(byteArray)

		assert.NoError(t, err)
		assert.Equal(t, DeviceId(0xffabcdef), deviceId)
	})

	t.Run(fmt.Sprintf("return an error when the array's length is grater than %d", SIZE_DEVICE_ID), func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd, 0xef, 0xaa}
		_, err := DeviceIdFromByteArray(byteArray)

		assert.Errorf(t, err, "invalid length (got %d, need: %d)", len(byteArray), SIZE_DEVICE_ID)
	})

	t.Run(fmt.Sprintf("return the padded DeviceId when the array's length is lesser than %d", SIZE_DEVICE_ID), func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd}
		deviceId, err := DeviceIdFromByteArray(byteArray)

		assert.NoError(t, err)
		assert.Equal(t, DeviceId(0xffabcd), deviceId)
	})
}
