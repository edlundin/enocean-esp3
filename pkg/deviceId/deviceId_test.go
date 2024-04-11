package deviceId

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceId_String(t *testing.T) {
	t.Run("returns lower case", func(t *testing.T) {
		id := DeviceId(0xffabcdef)
		assert.Equal(t, "ffabcdef", id.String())
	})
}

func TestDeviceId_DeviceIdFromHexString(t *testing.T) {
	id, err := DeviceIdFromHexString("ffabcdef")
	assert.NoError(t, err)
	assert.Equal(t, DeviceId(0xffabcdef), id)
}
