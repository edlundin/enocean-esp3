package enums

import "errors"

type MemoryType byte

const (
	MemoryTypeFLASH MemoryType = iota
	MemoryTypeRAM0
	MemoryTypeRAM_DATA
	MemoryTypeRAM_IDATA
	MemoryTypeRAM_XDATA
	MemoryTypeRAM_EEPROM
)

// ParseMemoryTypeFromByte parses a MemoryType from a byte.
func ParseMemoryTypeFromByte(b byte) (MemoryType, error) {
	switch b {
	case 0x00:
		return MemoryTypeFLASH, nil
	case 0x01:
		return MemoryTypeRAM0, nil
	case 0x02:
		return MemoryTypeRAM_DATA, nil
	case 0x03:
		return MemoryTypeRAM_IDATA, nil
	case 0x04:
		return MemoryTypeRAM_XDATA, nil
	case 0x05:
		return MemoryTypeRAM_EEPROM, nil
	default:
		return 0, errors.New("invalid memory type")
	}
}

// String returns the string representation of MemoryType.
func (memoryType MemoryType) String() string {
	switch memoryType {
	case MemoryTypeFLASH:
		return "FLASH"
	case MemoryTypeRAM0:
		return "RAM0"
	case MemoryTypeRAM_DATA:
		return "RAM_DATA"
	case MemoryTypeRAM_IDATA:
		return "RAM_IDATA"
	case MemoryTypeRAM_XDATA:
		return "RAM_XDATA"
	case MemoryTypeRAM_EEPROM:
		return "RAM_EEPROM"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether MemoryType is valid.
func (memoryType MemoryType) Valid() bool {
	switch memoryType {
	case MemoryTypeFLASH, MemoryTypeRAM0, MemoryTypeRAM_DATA, MemoryTypeRAM_IDATA, MemoryTypeRAM_XDATA, MemoryTypeRAM_EEPROM:
		return true
	default:
		return false
	}
}

type MemoryArea byte

const (
	MemoryAreaCONFIG MemoryArea = iota
	MemoryAreaSMART_ACK_TABLE
	MemoryAreaSYSTEM_ERROR_LOG
)

// ParseMemoryAreaFromByte parses a MemoryArea from a byte.
func ParseMemoryAreaFromByte(b byte) (MemoryArea, error) {
	switch b {
	case 0x00:
		return MemoryAreaCONFIG, nil
	case 0x01:
		return MemoryAreaSMART_ACK_TABLE, nil
	case 0x02:
		return MemoryAreaSYSTEM_ERROR_LOG, nil
	default:
		return 0, errors.New("invalid memory area")
	}
}

// String returns the string representation of MemoryArea.
func (memoryArea MemoryArea) String() string {
	switch memoryArea {
	case MemoryAreaCONFIG:
		return "CONFIG"
	case MemoryAreaSMART_ACK_TABLE:
		return "SMART_ACK_TABLE"
	case MemoryAreaSYSTEM_ERROR_LOG:
		return "SYSTEM_ERROR_LOG"
	default:
		return "UNKNOWN"
	}
}

// Valid reports whether MemoryArea is valid.
func (memoryArea MemoryArea) Valid() bool {
	switch memoryArea {
	case MemoryAreaCONFIG, MemoryAreaSMART_ACK_TABLE, MemoryAreaSYSTEM_ERROR_LOG:
		return true
	default:
		return false
	}
}
