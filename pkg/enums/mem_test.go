package enums

import "testing"

func TestMemoryType(t *testing.T) {
	cases := []struct{ b byte; v MemoryType; s string }{
		{0x00, MemoryTypeFLASH, "FLASH"}, {0x01, MemoryTypeRAM0, "RAM0"},
		{0x02, MemoryTypeRAM_DATA, "RAM_DATA"}, {0x03, MemoryTypeRAM_IDATA, "RAM_IDATA"},
		{0x04, MemoryTypeRAM_XDATA, "RAM_XDATA"}, {0x05, MemoryTypeRAM_EEPROM, "RAM_EEPROM"},
	}
	for _, c := range cases {
		v, err := ParseMemoryTypeFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseMemoryTypeFromByte(0xff); err == nil { t.Fatal("expected error") }
	if MemoryType(0xff).String() != "UNKNOWN" || MemoryType(0xff).Valid() { t.Fatal("invalid memory type accepted") }
}

func TestMemoryArea(t *testing.T) {
	cases := []struct{ b byte; v MemoryArea; s string }{
		{0x00, MemoryAreaCONFIG, "CONFIG"},
		{0x01, MemoryAreaSMART_ACK_TABLE, "SMART_ACK_TABLE"},
		{0x02, MemoryAreaSYSTEM_ERROR_LOG, "SYSTEM_ERROR_LOG"},
	}
	for _, c := range cases {
		v, err := ParseMemoryAreaFromByte(c.b)
		if err != nil || v != c.v || v.String() != c.s || !v.Valid() { t.Fatalf("%#x => %v %v", c.b, v, err) }
	}
	if _, err := ParseMemoryAreaFromByte(0xff); err == nil { t.Fatal("expected error") }
	if MemoryArea(0xff).String() != "UNKNOWN" || MemoryArea(0xff).Valid() { t.Fatal("invalid memory area accepted") }
}
