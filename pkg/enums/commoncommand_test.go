package enums

import "testing"

var commonCommandCases = []struct {
	value   byte
	command CommonCommand
	name    string
}{
	{0x01, CommonCommandWR_SLEEP, "WR_SLEEP"},
	{0x02, CommonCommandWR_RESET, "WR_RESET"},
	{0x03, CommonCommandRD_VERSION, "RD_VERSION"},
	{0x04, CommonCommandRD_SYS_LOG, "RD_SYS_LOG"},
	{0x05, CommonCommandRESET_SYS_LOG, "RESET_SYS_LOG"},
	{0x06, CommonCommandWR_BIST, "WR_BIST"},
	{0x07, CommonCommandWR_IDBASE, "WR_IDBASE"},
	{0x08, CommonCommandRD_IDBASE, "RD_IDBASE"},
	{0x09, CommonCommandWR_REPEATER, "WR_REPEATER"},
	{0x0a, CommonCommandRD_REPEATER, "RD_REPEATER"},
	{0x0b, CommonCommandWR_FILTER_ADD, "WR_FILTER_ADD"},
	{0x0c, CommonCommandWR_FILTER_DEL, "WR_FILTER_DEL"},
	{0x0d, CommonCommandWR_FILTER_DEL_ALL, "WR_FILTER_DEL_ALL"},
	{0x0e, CommonCommandWR_FILTER_ENABLE, "WR_FILTER_ENABLE"},
	{0x0f, CommonCommandRD_FILTER, "RD_FILTER"},
	{0x10, CommonCommandWR_WAIT_MATURITY, "WR_WAIT_MATURITY"},
	{0x11, CommonCommandWR_SUBTEL, "WR_SUBTEL"},
	{0x12, CommonCommandWR_MEM, "WR_MEM"},
	{0x13, CommonCommandRD_MEM, "RD_MEM"},
	{0x14, CommonCommandRD_MEM_ADDRESS, "RD_MEM_ADDRESS"},
	{0x15, CommonCommandRD_SECURITY, "RD_SECURITY"},
	{0x16, CommonCommandWR_SECURITY, "WR_SECURITY"},
	{0x17, CommonCommandWR_LEARNMODE, "WR_LEARNMODE"},
	{0x18, CommonCommandRD_LEARNMODE, "RD_LEARNMODE"},
	{0x19, CommonCommandWR_SECUREDEVICE_ADD, "WR_SECUREDEVICE_ADD"},
	{0x1a, CommonCommandWR_SECUREDEVICE_DEL, "WR_SECUREDEVICE_DEL"},
	{0x1b, CommonCommandRD_SECUREDEVICE_BY_INDEX, "RD_SECUREDEVICE_BY_INDEX"},
	{0x1c, CommonCommandWR_MODE, "WR_MODE"},
	{0x1d, CommonCommandRD_NUMSECUREDEVICES, "RD_NUMSECUREDEVICES"},
	{0x1e, CommonCommandRD_SECUREDEVICE_BY_ID, "RD_SECUREDEVICE_BY_ID"},
	{0x1f, CommonCommandWR_SECUREDEVICE_ADD_PSK, "WR_SECUREDEVICE_ADD_PSK"},
	{0x20, CommonCommandWR_SECUREDEVICE_SENDTEACHIN, "WR_SECUREDEVICE_SENDTEACHIN"},
	{0x21, CommonCommandWR_TEMPORARY_RLC_WINDOW, "WR_TEMPORARY_RLC_WINDOW"},
	{0x22, CommonCommandRD_SECUREDEVICE_PSK, "RD_SECUREDEVICE_PSK"},
	{0x23, CommonCommandRD_DUTYCYCLE_LIMIT, "RD_DUTYCYCLE_LIMIT"},
	{0x24, CommonCommandSET_BAUDRATE, "SET_BAUDRATE"},
	{0x25, CommonCommandGET_FREQUENCY_INFO, "GET_FREQUENCY_INFO"},
	{0x27, CommonCommandGET_STEPCODE, "GET_STEPCODE"},
	{0x2e, CommonCommandWR_REMAN_CODE, "WR_REMAN_CODE"},
	{0x2f, CommonCommandWR_STARTUP_DELAY, "WR_STARTUP_DELAY"},
	{0x30, CommonCommandWR_REMAN_REPEATING, "WR_REMAN_REPEATING"},
	{0x31, CommonCommandRD_REMAN_REPEATING, "RD_REMAN_REPEATING"},
	{0x32, CommonCommandSET_NOISETHRESHOLD, "SET_NOISETHRESHOLD"},
	{0x33, CommonCommandGET_NOISETHRESHOLD, "GET_NOISETHRESHOLD"},
	{0x34, CommonCommandSET_CRCMode, "SET_CRCMode"},
	{0x35, CommonCommandGET_CRCMode, "GET_CRCMode"},
	{0x36, CommonCommandWR_RLC_SAVE_PERIOD, "WR_RLC_SAVE_PERIOD"},
	{0x37, CommonCommandWR_RLC_LEGACY_MODE, "WR_RLC_LEGACY_MODE"},
	{0x38, CommonCommandWR_SECUREDEVICEV2_ADD, "WR_SECUREDEVICEV2_ADD"},
	{0x39, CommonCommandRD_SECUREDEVICEV2_BY_INDEX, "RD_SECUREDEVICEV2_BY_INDEX"},
	{0x3a, CommonCommandWR_RSSITEST_MODE, "WR_RSSITEST_MODE"},
	{0x3b, CommonCommandRD_RSSITEST_MODE, "RD_RSSITEST_MODE"},
	{0x3c, CommonCommandWR_SECUREDEVICE_REMAN_KEY, "WR_SECUREDEVICE_REMAN_KEY"},
	{0x3d, CommonCommandRD_SECUREDEVICE_REMAN_KEY, "RD_SECUREDEVICE_REMAN_KEY"},
	{0x3e, CommonCommandWR_TRANSPARENT_MODE, "WR_TRANSPARENT_MODE"},
	{0x3f, CommonCommandRD_TRANSPARENT_MODE, "RD_TRANSPARENT_MODE"},
	{0x40, CommonCommandWR_TX_ONLY_MODE, "WR_TX_ONLY_MODE"},
	{0x41, CommonCommandRD_TX_ONLY_MODE, "RD_TX_ONLY_MODE"},
}

func TestCommonCommands(t *testing.T) {
	for _, tc := range commonCommandCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseCommonCommandFromByte(tc.value)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.command || got.String() != tc.name || !got.Valid() {
				t.Fatalf("parsed=%v string=%q valid=%v", got, got.String(), got.Valid())
			}
		})
	}
}

func TestInvalidCommonCommands(t *testing.T) {
	valid := make(map[byte]bool, len(commonCommandCases))
	for _, tc := range commonCommandCases {
		valid[tc.value] = true
	}
	for i := 0; i <= 0xff; i++ {
		value := byte(i)
		if valid[value] {
			continue
		}
		got, err := ParseCommonCommandFromByte(value)
		if err == nil || got != 0 {
			t.Errorf("ParseCommonCommandFromByte(%02x) = %v, %v", value, got, err)
		}
		command := CommonCommand(value)
		if command.Valid() || command.String() != "UNKNOWN" {
			t.Errorf("CommonCommand(%02x): String=%q Valid=%v", value, command.String(), command.Valid())
		}
	}
}

type modeValue interface {
	comparable
	String() string
	Valid() bool
}

func TestCommonCommandModes(t *testing.T) {
	t.Run("baud", func(t *testing.T) {
		testMode(t, []TCMBaudrate{TCMBaudrate57600, TCMBaudrate115200, TCMBaudrate230400, TCMBaudrate460800}, []string{"57600", "115200", "230400", "460800"}, TCMBaudrate(0xff), ParseTCMBaudrateFromByte)
	})
	t.Run("frequency", func(t *testing.T) {
		testMode(t, []TCMFrequency{TCMFrequency315_000_MHZ, TCMFrequency868_000_MHZ, TCMFrequency902_875_MHZ, TCMFrequency921_400_MHZ, TCMFrequency928_350_MHZ, TCMFrequency2_4_GHZ}, []string{"315.000 MHz", "868.000 MHz", "902.875 MHz", "921.400 MHz", "928.350 MHz", "2.4 GHz"}, TCMFrequency(0xff), nil)
	})
	t.Run("protocol", func(t *testing.T) {
		testMode(t, []TCMProtocol{TCMProtocolERP1, TCMProtocolERP2, TCMProtocolIEEE_802_15_4, TCMProtocolLONG_RANGE}, []string{"ERP1", "ERP2", "IEEE 802.15.4", "LONG_RANGE"}, TCMProtocol(0xff), nil)
	})
	t.Run("CRC", func(t *testing.T) {
		testMode(t, []CRCMode{CRCMode8BIT, CRCMode7BIT}, []string{"8BIT", "7BIT"}, CRCMode(0xff), ParseCRCModeFromByte)
	})
	t.Run("RLC", func(t *testing.T) {
		testMode(t, []RLCMode{RLCModeSTANDARD, RLCModeLEGACY}, []string{"STANDARD", "LEGACY"}, RLCMode(0xff), ParseRLCModeFromByte)
	})
	t.Run("RSSI", func(t *testing.T) {
		testMode(t, []RSSITestMode{RSSITestModeDISABLED, RSSITestModeENABLED}, []string{"DISABLED", "ENABLED"}, RSSITestMode(0xff), ParseRSSITestModeFromByte)
	})
	t.Run("transparent", func(t *testing.T) {
		testMode(t, []TransparentMode{TransparentModeDISABLED, TransparentModeENABLED}, []string{"DISABLED", "ENABLED"}, TransparentMode(0xff), ParseTransparentModeFromByte)
	})
	t.Run("TX-only", func(t *testing.T) {
		testMode(t, []TxOnlyMode{TxOnlyModeDISABLED, TxOnlyModeENABLED_WITHOUT_AUTO_SLEEP, TxOnlyModeENABLED_WITH_AUTO_SLEEP}, []string{"DISABLED", "ENABLED_WITHOUT_AUTO_SLEEP", "ENABLED_WITH_AUTO_SLEEP"}, TxOnlyMode(0xff), ParseTxOnlyModeFromByte)
	})
}

func testMode[T modeValue](t *testing.T, values []T, names []string, invalid T, parse func(byte) (T, error)) {
	t.Helper()
	for i, value := range values {
		if got := value.String(); got != names[i] {
			t.Errorf("value %d: String() = %q, want %q", i, got, names[i])
		}
		if !value.Valid() {
			t.Errorf("value %d is invalid", i)
		}
		if parse != nil {
			got, err := parse(byte(i))
			if err != nil || got != value {
				t.Errorf("parse(%02x) = %v, %v; want %v", i, got, err, value)
			}
		}
	}
	if invalid.String() != "UNKNOWN" || invalid.Valid() {
		t.Errorf("invalid value: String=%q Valid=%v", invalid.String(), invalid.Valid())
	}
	if parse != nil {
		if _, err := parse(0xff); err == nil {
			t.Error("parse(ff) succeeded")
		}
	}
}
