package enums

import (
	"testing"
)

func TestReturnCodeConstants(t *testing.T) {
	tests := []struct {
		name     string
		code     ReturnCode
		expected byte
	}{
		{"SUCCESS", ReturnCodeSUCCESS, 0x00},
		{"ERROR", ReturnCodeERROR, 0x01},
		{"NOT_SUPPORTED", ReturnCodeNOT_SUPPORTED, 0x02},
		{"WRONG_ARGUMENT", ReturnCodeWRONG_ARGUMENT, 0x03},
		{"OPERATION_DENIED", ReturnCodeOPERATION_DENIED, 0x04},
		{"LOCK_SET", ReturnCodeLOCK_SET, 0x05},
		{"BUFFER_TO_SMALL", ReturnCodeBUFFER_TO_SMALL, 0x06},
		{"NO_FREE_BUFFER", ReturnCodeNO_FREE_BUFFER, 0x07},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if byte(tt.code) != tt.expected {
				t.Errorf("ReturnCode%s = %d, want %d", tt.name, byte(tt.code), tt.expected)
			}
		})
	}
}

func TestParseReturnCodeFromByte(t *testing.T) {
	tests := []struct {
		name        string
		input       uint8
		expected    ReturnCode
		expectError bool
	}{
		{"SUCCESS", 0x00, ReturnCodeSUCCESS, false},
		{"ERROR", 0x01, ReturnCodeERROR, false},
		{"NOT_SUPPORTED", 0x02, ReturnCodeNOT_SUPPORTED, false},
		{"WRONG_ARGUMENT", 0x03, ReturnCodeWRONG_ARGUMENT, false},
		{"OPERATION_DENIED", 0x04, ReturnCodeOPERATION_DENIED, false},
		{"LOCK_SET", 0x05, ReturnCodeLOCK_SET, false},
		{"BUFFER_TO_SMALL", 0x06, ReturnCodeBUFFER_TO_SMALL, false},
		{"NO_FREE_BUFFER", 0x07, ReturnCodeNO_FREE_BUFFER, false},
		{"Invalid code 0x08", 0x08, ReturnCodeERROR, true},
		{"Invalid code 0xFF", 0xFF, ReturnCodeERROR, true},
		{"Invalid code 0x10", 0x10, ReturnCodeERROR, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseReturnCodeFromByte(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParseReturnCodeFromByte(%d) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseReturnCodeFromByte(%d) unexpected error: %v", tt.input, err)
				}
			}

			if result != tt.expected {
				t.Errorf("ParseReturnCodeFromByte(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestReturnCodeString(t *testing.T) {
	tests := []struct {
		name     string
		code     ReturnCode
		expected string
	}{
		{"SUCCESS", ReturnCodeSUCCESS, "SUCCESS"},
		{"ERROR", ReturnCodeERROR, "ERROR"},
		{"NOT_SUPPORTED", ReturnCodeNOT_SUPPORTED, "NOT_SUPPORTED"},
		{"WRONG_ARGUMENT", ReturnCodeWRONG_ARGUMENT, "WRONG_ARGUMENT"},
		{"OPERATION_DENIED", ReturnCodeOPERATION_DENIED, "OPERATION_DENIED"},
		{"LOCK_SET", ReturnCodeLOCK_SET, "LOCK_SET"},
		{"BUFFER_TO_SMALL", ReturnCodeBUFFER_TO_SMALL, "BUFFER_TO_SMALL"},
		{"NO_FREE_BUFFER", ReturnCodeNO_FREE_BUFFER, "NO_FREE_BUFFER"},
		{"Unknown code", ReturnCode(0xFF), "UNKNOWN"},
		{"Unknown code 0x08", ReturnCode(0x08), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.code.String()
			if result != tt.expected {
				t.Errorf("ReturnCode(%d).String() = %s, want %s", byte(tt.code), result, tt.expected)
			}
		})
	}
}

func TestReturnCodeRoundTrip(t *testing.T) {
	// Test that parsing a byte and converting back gives the same result
	codes := []ReturnCode{
		ReturnCodeSUCCESS,
		ReturnCodeERROR,
		ReturnCodeNOT_SUPPORTED,
		ReturnCodeWRONG_ARGUMENT,
		ReturnCodeOPERATION_DENIED,
		ReturnCodeLOCK_SET,
		ReturnCodeBUFFER_TO_SMALL,
		ReturnCodeNO_FREE_BUFFER,
	}

	for _, code := range codes {
		t.Run(code.String(), func(t *testing.T) {
			parsed, err := ParseReturnCodeFromByte(uint8(code))
			if err != nil {
				t.Errorf("ParseReturnCodeFromByte(%d) failed: %v", byte(code), err)
			}
			if parsed != code {
				t.Errorf("Round trip failed: %d -> %d", byte(code), byte(parsed))
			}
		})
	}
}

func TestReturnCodeStringRoundTrip(t *testing.T) {
	// Test that converting to string and back gives consistent results
	codes := []ReturnCode{
		ReturnCodeSUCCESS,
		ReturnCodeERROR,
		ReturnCodeNOT_SUPPORTED,
		ReturnCodeWRONG_ARGUMENT,
		ReturnCodeOPERATION_DENIED,
		ReturnCodeLOCK_SET,
		ReturnCodeBUFFER_TO_SMALL,
		ReturnCodeNO_FREE_BUFFER,
	}

	for _, code := range codes {
		t.Run(code.String(), func(t *testing.T) {
			str := code.String()
			if str == "" {
				t.Errorf("String() returned empty string for code %d", byte(code))
			}
			if str == "UNKNOWN" {
				t.Errorf("String() returned UNKNOWN for valid code %d", byte(code))
			}
		})
	}
}
