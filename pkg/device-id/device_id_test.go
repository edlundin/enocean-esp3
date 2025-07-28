package device_id

import (
	"fmt"
	"testing"
)

func TestDeviceId_GetBroadcastId(t *testing.T) {
	t.Run("returns 0xffffffff (broadcast ID)", func(t *testing.T) {
		if BroadcastId() != DeviceID(0xffffffff) {
			t.Errorf("expected %d, got %d", 0xffffffff, BroadcastId())
		}
	})
}

func TestDeviceId_String(t *testing.T) {
	t.Run("returns lower case without hex prefix", func(t *testing.T) {
		deviceId := DeviceID(0xffabc123)

		if deviceId.String() != "ffabc123" {
			t.Errorf("expected: %s, got: %s", "ffabc123", deviceId.String())
		}
	})

	t.Run("always returns 8 characters, padded with 0", func(t *testing.T) {
		if deviceId := DeviceID(0); deviceId.String() != "00000000" {
			t.Errorf("expected: %s, got: %s", "00000000", deviceId.String())
		}

		if deviceId := DeviceID(0x42); deviceId.String() != "00000042" {
			t.Errorf("expected: %s, got: %s", "00000042", deviceId.String())
		}

		if deviceId := DeviceID(0xabcdef42); deviceId.String() != "abcdef42" {
			t.Errorf("expected: %s, got: %s", "abcdef42", deviceId.String())
		}
	})
}

func TestDeviceIdFromHexString(t *testing.T) {
	const sizeMaxStr = DeviceIDSize * 2

	t.Run("returns the DeviceID representation of the string content (even length)", func(t *testing.T) {
		deviceId, err := FromHexString("ffabcdef")

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if deviceId != DeviceID(0xffabcdef) {
			t.Errorf("expected %d, got %d", 0xffabcdef, deviceId)
		}
	})

	t.Run("return the DeviceID representation of the string content (odd length with padding)", func(t *testing.T) {
		deviceId, err := FromHexString("0xffa")

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if deviceId != DeviceID(0xffa) {
			t.Errorf("expected %d, got %d", 0xffa, deviceId)
		}
	})

	t.Run("ignore 0x prefix if present", func(t *testing.T) {
		idStr := "0xffabcdef"
		deviceId, err := FromHexString(idStr)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if deviceId != DeviceID(0xffabcdef) {
			t.Errorf("expected %d, got %d", 0xffabcdef, deviceId)
		}
	})

	t.Run(fmt.Sprintf("return an error when the string's length is greater than %d", sizeMaxStr), func(t *testing.T) {
		idStr := "ffabcdefaa"
		expectedError := fmt.Sprintf("invalid length (got:%d, max:%d)", len(idStr), sizeMaxStr)
		_, err := FromHexString(idStr)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if err.Error() != expectedError {
			t.Errorf("expected: %s, got: %s", expectedError, err.Error())
		}
	})

	t.Run("return an error when the string is invalid (not in hexadecimal format)", func(t *testing.T) {
		idStr := "ffabcdeg"
		expectedError := "encoding/hex: invalid byte: U+0067 'g'"
		_, err := FromHexString(idStr)

		if err == nil {
			t.Errorf("expected error, got nil")

		}

		if err.Error() != expectedError {
			t.Errorf("expected: %s, got: %s", expectedError, err.Error())
		}
	})
}

func TestDeviceIdFromByteArray(t *testing.T) {
	t.Run("returns the DeviceID representation of the array content", func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd, 0xef}
		deviceId, err := FromByteArray(byteArray)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if deviceId != DeviceID(0xffabcdef) {
			t.Errorf("expected %d, got %d", 0xffabcdef, deviceId)
		}
	})

	t.Run(fmt.Sprintf("return an error when the array's length is greater than %d", DeviceIDSize), func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd, 0xef, 0xaa}
		expectedError := fmt.Sprintf("invalid length (got:%d, need:%d)", len(byteArray), DeviceIDSize)
		_, err := FromByteArray(byteArray)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if err.Error() != expectedError {
			t.Errorf("expected: %s, got: %s", expectedError, err.Error())
		}
	})

	t.Run(fmt.Sprintf("return the padded DeviceID when the array's length is lesser than %d", DeviceIDSize), func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd}
		deviceId, err := FromByteArray(byteArray)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if deviceId != DeviceID(0xffabcd) {
			t.Errorf("expected %d, got %d", 0xffabcd, deviceId)
		}
	})
}

func TestDeviceId_ToArray(t *testing.T) {
	t.Run("returns big-endian byte array representation", func(t *testing.T) {
		deviceId := DeviceID(0xffabcdef)
		expected := [4]byte{0xff, 0xab, 0xcd, 0xef}
		result := deviceId.ToArray()

		if result != expected {
			t.Errorf("expected: %v, got: %v", expected, result)
		}
	})

	t.Run("returns correct array for zero value", func(t *testing.T) {
		deviceId := DeviceID(0)
		expected := [4]byte{0x00, 0x00, 0x00, 0x00}
		result := deviceId.ToArray()

		if result != expected {
			t.Errorf("expected: %v, got: %v", expected, result)
		}
	})

	t.Run("returns correct array for maximum value", func(t *testing.T) {
		deviceId := DeviceID(0xffffffff)
		expected := [4]byte{0xff, 0xff, 0xff, 0xff}
		result := deviceId.ToArray()

		if result != expected {
			t.Errorf("expected: %v, got: %v", expected, result)
		}
	})

	t.Run("returns correct array for single byte value", func(t *testing.T) {
		deviceId := DeviceID(0x42)
		expected := [4]byte{0x00, 0x00, 0x00, 0x42}
		result := deviceId.ToArray()

		if result != expected {
			t.Errorf("expected: %v, got: %v", expected, result)
		}
	})

	t.Run("returns correct array for two byte value", func(t *testing.T) {
		deviceId := DeviceID(0xabcd)
		expected := [4]byte{0x00, 0x00, 0xab, 0xcd}
		result := deviceId.ToArray()

		if result != expected {
			t.Errorf("expected: %v, got: %v", expected, result)
		}
	})

	t.Run("returns correct array for three byte value", func(t *testing.T) {
		deviceId := DeviceID(0xabcdef)
		expected := [4]byte{0x00, 0xab, 0xcd, 0xef}
		result := deviceId.ToArray()

		if result != expected {
			t.Errorf("expected: %v, got: %v", expected, result)
		}
	})

	t.Run("returns correct array for broadcast ID", func(t *testing.T) {
		deviceId := BroadcastId()
		expected := [4]byte{0xff, 0xff, 0xff, 0xff}
		result := deviceId.ToArray()

		if result != expected {
			t.Errorf("expected: %v, got: %v", expected, result)
		}
	})
}
