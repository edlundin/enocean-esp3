package enums

import "errors"

type SecureDeviceDirection byte

const (
	SecureDeviceDirectionINBOUND_TABLE            SecureDeviceDirection = 0x00
	SecureDeviceDirectionOUTBOUND_TABLE           SecureDeviceDirection = 0x01
	SecureDeviceDirectionOUTBOUND_BROADCAST_TABLE SecureDeviceDirection = 0x02
	// NOTE: Only used for RdNumSecureDevices
	SecureDeviceDirectionALL SecureDeviceDirection = 0x03
	// NOTE: Only used for RdSecureDeviceByID
	SecureDeviceDirectionREMAN_TABLE SecureDeviceDirection = 0x03
	SecureDeviceDirectionNONE        SecureDeviceDirection = 0xff
)

// ParseSecureDeviceDirectionFromByte parses a SecureDeviceDirection from a byte.
func ParseSecureDeviceDirectionFromByte(b byte) (SecureDeviceDirection, error) {
	switch b {
	case 0x00:
		return SecureDeviceDirectionINBOUND_TABLE, nil
	case 0x01:
		return SecureDeviceDirectionOUTBOUND_TABLE, nil
	case 0x02:
		return SecureDeviceDirectionOUTBOUND_BROADCAST_TABLE, nil
	case 0x03:
		return SecureDeviceDirectionALL, nil
	case 0xff:
		return SecureDeviceDirectionNONE, nil
	default:
		return 0, errors.New("invalid secure device direction")
	}
}

// String returns the string representation of SecureDeviceDirection.
func (direction SecureDeviceDirection) String() string {
	switch direction {
	case SecureDeviceDirectionINBOUND_TABLE:
		return "INBOUND_TABLE"
	case SecureDeviceDirectionOUTBOUND_TABLE:
		return "OUTBOUND_TABLE"
	case SecureDeviceDirectionOUTBOUND_BROADCAST_TABLE:
		return "OUTBOUND_BROADCAST_TABLE"
	case SecureDeviceDirectionALL:
		return "ALL_OR_REMAN_TABLE"
	case SecureDeviceDirectionNONE:
		return "NONE"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether SecureDeviceDirection is valid.
func (direction SecureDeviceDirection) Valid() bool {
	switch direction {
	case SecureDeviceDirectionINBOUND_TABLE,
		SecureDeviceDirectionOUTBOUND_TABLE,
		SecureDeviceDirectionOUTBOUND_BROADCAST_TABLE,
		SecureDeviceDirectionALL,
		SecureDeviceDirectionNONE:
		return true
	default:
		return false
	}
}
