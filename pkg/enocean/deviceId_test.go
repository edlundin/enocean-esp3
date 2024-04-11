package enocean

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceId_String(t *testing.T) {
	t.Run("returns lower case without prefix", func(t *testing.T) {
		deviceId := DeviceId(0xffabc123)
		assert.Equal(t, "ffabc123", deviceId.String())
	})

	t.Run("always returns 8 characters, padded with 0", func(t *testing.T) {
		deviceId := DeviceId(0)
		assert.Equal(t, "00000000", deviceId.String())

		deviceId = DeviceId(0x42)
		assert.Equal(t, "00000042", deviceId.String())
	})
}

func TestDeviceId_DeviceIdFromHexString(t *testing.T) {
	id, err := DeviceIdFromHexString("ffabcdef")
	assert.NoError(t, err)
	assert.Equal(t, DeviceId(0xffabcdef), id)
}
