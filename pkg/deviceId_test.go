package pkg

import (
	"fmt"
	"testing"
)

func TestDeviceId_GetBroadcastId(t *testing.T) {
	t.Run("returns 0xffffffff (broadcast ID)", func(t *testing.T) {
		if GetBroadcastId() != DeviceID(0xffffffff) {
			t.Errorf("expected %d, got %d", 0xffffffff, GetBroadcastId())
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
	const sizeMaxStr = sizeDeviceID * 2

	t.Run("returns the DeviceID representation of the string content (even length)", func(t *testing.T) {
		deviceId, err := DeviceIdFromHexString("ffabcdef")

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if deviceId != DeviceID(0xffabcdef) {
			t.Errorf("expected %d, got %d", 0xffabcdef, deviceId)
		}
	})

	t.Run("return the DeviceID representation of the string content (odd length with padding)", func(t *testing.T) {
		deviceId, err := DeviceIdFromHexString("0xffa")

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if deviceId != DeviceID(0xffa) {
			t.Errorf("expected %d, got %d", 0xffa, deviceId)
		}
	})

	t.Run("ignore 0x prefix if present", func(t *testing.T) {
		idStr := "0xffabcdef"
		deviceId, err := DeviceIdFromHexString(idStr)

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
		_, err := DeviceIdFromHexString(idStr)

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
		_, err := DeviceIdFromHexString(idStr)

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
		deviceId, err := DeviceIdFromByteArray(byteArray)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if deviceId != DeviceID(0xffabcdef) {
			t.Errorf("expected %d, got %d", 0xffabcdef, deviceId)
		}
	})

	t.Run(fmt.Sprintf("return an error when the array's length is greater than %d", sizeDeviceID), func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd, 0xef, 0xaa}
		expectedError := fmt.Sprintf("invalid length (got:%d, need:%d)", len(byteArray), sizeDeviceID)
		_, err := DeviceIdFromByteArray(byteArray)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if err.Error() != expectedError {
			t.Errorf("expected: %s, got: %s", expectedError, err.Error())
		}
	})

	t.Run(fmt.Sprintf("return the padded DeviceID when the array's length is lesser than %d", sizeDeviceID), func(t *testing.T) {
		byteArray := []byte{0xff, 0xab, 0xcd}
		deviceId, err := DeviceIdFromByteArray(byteArray)

		if err != nil {
			t.Errorf("expected no error, got: %s", err)
		}

		if deviceId != DeviceID(0xffabcd) {
			t.Errorf("expected %d, got %d", 0xffabcd, deviceId)
		}
	})
}
