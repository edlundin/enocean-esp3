package enums

import "errors"

type CommonCommand byte

const (
	CommonCommandWR_SLEEP                       CommonCommand = 0x01
	CommonCommandWR_RESET                       CommonCommand = 0x02
	CommonCommandRD_VERSION                     CommonCommand = 0x03
	CommonCommandRD_SYS_LOG                     CommonCommand = 0x04
	CommonCommandRESET_SYS_LOG                  CommonCommand = 0x05
	CommonCommandWR_BIST                        CommonCommand = 0x06
	CommonCommandWR_IDBASE                      CommonCommand = 0x07
	CommonCommandRD_IDBASE                      CommonCommand = 0x08
	CommonCommandWR_REPEATER                    CommonCommand = 0x09
	CommonCommandRD_REPEATER                    CommonCommand = 0x0a
	CommonCommandWR_FILTER_ADD                  CommonCommand = 0x0b
	CommonCommandWR_FILTER_DEL                  CommonCommand = 0x0c
	CommonCommandWR_FILTER_DEL_ALL              CommonCommand = 0x0d
	CommonCommandWR_FILTER_ENABLE               CommonCommand = 0x0e
	CommonCommandRD_FILTER                      CommonCommand = 0x0f
	CommonCommandWR_WAIT_MATURITY               CommonCommand = 0x10
	CommonCommandWR_SUBTEL                      CommonCommand = 0x11
	CommonCommandWR_MEM                         CommonCommand = 0x12
	CommonCommandRD_MEM                         CommonCommand = 0x13
	CommonCommandRD_MEM_ADDRESS                 CommonCommand = 0x14
	CommonCommandRD_SECURITY                    CommonCommand = 0x15
	CommonCommandWR_SECURITY                    CommonCommand = 0x16
	CommonCommandWR_LEARNMODE                   CommonCommand = 0x17
	CommonCommandRD_LEARNMODE                   CommonCommand = 0x18
	CommonCommandWR_SECUREDEVICE_ADD            CommonCommand = 0x19
	CommonCommandWR_SECUREDEVICE_DEL            CommonCommand = 0x1a
	CommonCommandRD_SECUREDEVICE_BY_INDEX       CommonCommand = 0x1b
	CommonCommandWR_MODE                        CommonCommand = 0x1c
	CommonCommandRD_NUMSECUREDEVICES            CommonCommand = 0x1d
	CommonCommandRD_SECUREDEVICE_BY_ID          CommonCommand = 0x1e
	CommonCommandWR_SECUREDEVICE_ADD_PSK        CommonCommand = 0x1f
	CommonCommandWR_SECUREDEVICE_SENDTEACHIN    CommonCommand = 0x20
	CommonCommandWR_TEMPORARY_RLC_WINDOW        CommonCommand = 0x21
	CommonCommandRD_SECUREDEVICE_PSK            CommonCommand = 0x22
	CommonCommandRD_DUTYCYCLE_LIMIT             CommonCommand = 0x23
	CommonCommandSET_BAUDRATE                   CommonCommand = 0x24
	CommonCommandGET_FREQUENCY_INFO             CommonCommand = 0x25
	CommonCommandGET_STEPCODE                   CommonCommand = 0x27
	CommonCommandWR_REMAN_CODE                  CommonCommand = 0x2e
	CommonCommandWR_STARTUP_DELAY               CommonCommand = 0x2f
	CommonCommandWR_REMAN_REPEATING             CommonCommand = 0x30
	CommonCommandRD_REMAN_REPEATING             CommonCommand = 0x31
	CommonCommandSET_NOISETHRESHOLD             CommonCommand = 0x32
	CommonCommandGET_NOISETHRESHOLD             CommonCommand = 0x33
	CommonCommandWR_RLC_SAVE_PERIOD             CommonCommand = 0x36
	CommonCommandWR_RLC_LEGACY_MODE             CommonCommand = 0x37
	CommonCommandWR_SECUREDEVICEV2_ADD          CommonCommand = 0x38
	CommonCommandRD_SECUREDEVICEV2_BY_INDEX     CommonCommand = 0x39
	CommonCommandWR_RSSITEST_MODE               CommonCommand = 0x3a
	CommonCommandRD_RSSITEST_MODE               CommonCommand = 0x3b
	CommonCommandWR_SECUREDEVICE_MAINTENANCEKEY CommonCommand = 0x3c
	CommonCommandRD_SECUREDEVICE_MAINTENANCEKEY CommonCommand = 0x3d
	CommonCommandWR_TRANSPARENT_MODE            CommonCommand = 0x3e
	CommonCommandRD_TRANSPARENT_MODE            CommonCommand = 0x3f
	CommonCommandWR_TX_ONLY_MODE                CommonCommand = 0x40
	CommonCommandRD_TX_ONLY_MODE                CommonCommand = 0x41
)

func ParseCommonCommandFromByte(b byte) (CommonCommand, error) {
	switch b {
	case 0x01:
		return CommonCommandWR_SLEEP, nil
	case 0x02:
		return CommonCommandWR_RESET, nil
	case 0x03:
		return CommonCommandRD_VERSION, nil
	case 0x04:
		return CommonCommandRD_SYS_LOG, nil
	case 0x05:
		return CommonCommandRESET_SYS_LOG, nil
	case 0x06:
		return CommonCommandWR_BIST, nil
	case 0x07:
		return CommonCommandWR_IDBASE, nil
	case 0x08:
		return CommonCommandRD_IDBASE, nil
	case 0x09:
		return CommonCommandWR_REPEATER, nil
	case 0x0a:
		return CommonCommandRD_REPEATER, nil
	case 0x0b:
		return CommonCommandWR_FILTER_ADD, nil
	case 0x0c:
		return CommonCommandWR_FILTER_DEL, nil
	case 0x0d:
		return CommonCommandWR_FILTER_DEL_ALL, nil
	case 0x0e:
		return CommonCommandWR_FILTER_ENABLE, nil
	case 0x0f:
		return CommonCommandRD_FILTER, nil
	case 0x10:
		return CommonCommandWR_WAIT_MATURITY, nil
	case 0x11:
		return CommonCommandWR_SUBTEL, nil
	case 0x12:
		return CommonCommandWR_MEM, nil
	case 0x13:
		return CommonCommandRD_MEM, nil
	case 0x14:
		return CommonCommandRD_MEM_ADDRESS, nil
	case 0x15:
		return CommonCommandRD_SECURITY, nil
	case 0x16:
		return CommonCommandWR_SECURITY, nil
	case 0x17:
		return CommonCommandWR_LEARNMODE, nil
	case 0x18:
		return CommonCommandRD_LEARNMODE, nil
	case 0x19:
		return CommonCommandWR_SECUREDEVICE_ADD, nil
	case 0x1a:
		return CommonCommandWR_SECUREDEVICE_DEL, nil
	case 0x1b:
		return CommonCommandRD_SECUREDEVICE_BY_INDEX, nil
	case 0x1c:
		return CommonCommandWR_MODE, nil
	case 0x1d:
		return CommonCommandRD_NUMSECUREDEVICES, nil
	case 0x1e:
		return CommonCommandRD_SECUREDEVICE_BY_ID, nil
	case 0x1f:
		return CommonCommandWR_SECUREDEVICE_ADD_PSK, nil
	case 0x20:
		return CommonCommandWR_SECUREDEVICE_SENDTEACHIN, nil
	case 0x21:
		return CommonCommandWR_TEMPORARY_RLC_WINDOW, nil
	case 0x22:
		return CommonCommandRD_SECUREDEVICE_PSK, nil
	case 0x23:
		return CommonCommandRD_DUTYCYCLE_LIMIT, nil
	case 0x24:
		return CommonCommandSET_BAUDRATE, nil
	case 0x25:
		return CommonCommandGET_FREQUENCY_INFO, nil
	case 0x27:
		return CommonCommandGET_STEPCODE, nil
	case 0x2e:
		return CommonCommandWR_REMAN_CODE, nil
	case 0x2f:
		return CommonCommandWR_STARTUP_DELAY, nil
	case 0x30:
		return CommonCommandWR_REMAN_REPEATING, nil
	case 0x31:
		return CommonCommandRD_REMAN_REPEATING, nil
	case 0x32:
		return CommonCommandSET_NOISETHRESHOLD, nil
	case 0x33:
		return CommonCommandGET_NOISETHRESHOLD, nil
	case 0x36:
		return CommonCommandWR_RLC_SAVE_PERIOD, nil
	case 0x37:
		return CommonCommandWR_RLC_LEGACY_MODE, nil
	case 0x38:
		return CommonCommandWR_SECUREDEVICEV2_ADD, nil
	case 0x39:
		return CommonCommandRD_SECUREDEVICEV2_BY_INDEX, nil
	case 0x3a:
		return CommonCommandWR_RSSITEST_MODE, nil
	case 0x3b:
		return CommonCommandRD_RSSITEST_MODE, nil
	case 0x3c:
		return CommonCommandWR_SECUREDEVICE_MAINTENANCEKEY, nil
	case 0x3d:
		return CommonCommandRD_SECUREDEVICE_MAINTENANCEKEY, nil
	case 0x3e:
		return CommonCommandWR_TRANSPARENT_MODE, nil
	case 0x3f:
		return CommonCommandRD_TRANSPARENT_MODE, nil
	case 0x40:
		return CommonCommandWR_TX_ONLY_MODE, nil
	case 0x41:
		return CommonCommandRD_TX_ONLY_MODE, nil
	default:
		return 0, errors.New("invalid common command")
	}
}

func (command CommonCommand) String() string {
	switch command {
	case CommonCommandWR_SLEEP:
		return "WR_SLEEP"
	case CommonCommandWR_RESET:
		return "WR_RESET"
	case CommonCommandRD_VERSION:
		return "RD_VERSION"
	case CommonCommandRD_SYS_LOG:
		return "RD_SYS_LOG"
	case CommonCommandRESET_SYS_LOG:
		return "RESET_SYS_LOG"
	case CommonCommandWR_BIST:
		return "WR_BIST"
	case CommonCommandWR_IDBASE:
		return "WR_IDBASE"
	case CommonCommandRD_IDBASE:
		return "RD_IDBASE"
	case CommonCommandWR_REPEATER:
		return "WR_REPEATER"
	case CommonCommandRD_REPEATER:
		return "RD_REPEATER"
	case CommonCommandWR_FILTER_ADD:
		return "WR_FILTER_ADD"
	case CommonCommandWR_FILTER_DEL:
		return "WR_FILTER_DEL"
	case CommonCommandWR_FILTER_DEL_ALL:
		return "WR_FILTER_DEL_ALL"
	case CommonCommandWR_FILTER_ENABLE:
		return "WR_FILTER_ENABLE"
	case CommonCommandRD_FILTER:
		return "RD_FILTER"
	case CommonCommandWR_WAIT_MATURITY:
		return "WR_WAIT_MATURITY"
	case CommonCommandWR_SUBTEL:
		return "WR_SUBTEL"
	case CommonCommandWR_MEM:
		return "WR_MEM"
	case CommonCommandRD_MEM:
		return "RD_MEM"
	case CommonCommandRD_MEM_ADDRESS:
		return "RD_MEM_ADDRESS"
	case CommonCommandRD_SECURITY:
		return "RD_SECURITY"
	case CommonCommandWR_SECURITY:
		return "WR_SECURITY"
	case CommonCommandWR_LEARNMODE:
		return "WR_LEARNMODE"
	case CommonCommandRD_LEARNMODE:
		return "RD_LEARNMODE"
	case CommonCommandWR_SECUREDEVICE_ADD:
		return "WR_SECUREDEVICE_ADD"
	case CommonCommandWR_SECUREDEVICE_DEL:
		return "WR_SECUREDEVICE_DEL"
	case CommonCommandRD_SECUREDEVICE_BY_INDEX:
		return "RD_SECUREDEVICE_BY_INDEX"
	case CommonCommandWR_MODE:
		return "WR_MODE"
	case CommonCommandRD_NUMSECUREDEVICES:
		return "RD_NUMSECUREDEVICES"
	case CommonCommandRD_SECUREDEVICE_BY_ID:
		return "RD_SECUREDEVICE_BY_ID"
	case CommonCommandWR_SECUREDEVICE_ADD_PSK:
		return "WR_SECUREDEVICE_ADD_PSK"
	case CommonCommandWR_SECUREDEVICE_SENDTEACHIN:
		return "WR_SECUREDEVICE_SENDTEACHIN"
	case CommonCommandWR_TEMPORARY_RLC_WINDOW:
		return "WR_TEMPORARY_RLC_WINDOW"
	case CommonCommandRD_SECUREDEVICE_PSK:
		return "RD_SECUREDEVICE_PSK"
	case CommonCommandRD_DUTYCYCLE_LIMIT:
		return "RD_DUTYCYCLE_LIMIT"
	case CommonCommandSET_BAUDRATE:
		return "SET_BAUDRATE"
	case CommonCommandGET_FREQUENCY_INFO:
		return "GET_FREQUENCY_INFO"
	case CommonCommandGET_STEPCODE:
		return "GET_STEPCODE"
	case CommonCommandWR_REMAN_CODE:
		return "WR_REMAN_CODE"
	case CommonCommandWR_STARTUP_DELAY:
		return "WR_STARTUP_DELAY"
	case CommonCommandWR_REMAN_REPEATING:
		return "WR_REMAN_REPEATING"
	case CommonCommandRD_REMAN_REPEATING:
		return "RD_REMAN_REPEATING"
	case CommonCommandSET_NOISETHRESHOLD:
		return "SET_NOISETHRESHOLD"
	case CommonCommandGET_NOISETHRESHOLD:
		return "GET_NOISETHRESHOLD"
	case CommonCommandWR_RLC_SAVE_PERIOD:
		return "WR_RLC_SAVE_PERIOD"
	case CommonCommandWR_RLC_LEGACY_MODE:
		return "WR_RLC_LEGACY_MODE"
	case CommonCommandWR_SECUREDEVICEV2_ADD:
		return "WR_SECUREDEVICEV2_ADD"
	case CommonCommandRD_SECUREDEVICEV2_BY_INDEX:
		return "RD_SECUREDEVICEV2_BY_INDEX"
	case CommonCommandWR_RSSITEST_MODE:
		return "WR_RSSITEST_MODE"
	case CommonCommandRD_RSSITEST_MODE:
		return "RD_RSSITEST_MODE"
	case CommonCommandWR_SECUREDEVICE_MAINTENANCEKEY:
		return "WR_SECUREDEVICE_MAINTENANCEKEY"
	case CommonCommandRD_SECUREDEVICE_MAINTENANCEKEY:
		return "RD_SECUREDEVICE_MAINTENANCEKEY"
	case CommonCommandWR_TRANSPARENT_MODE:
		return "WR_TRANSPARENT_MODE"
	case CommonCommandRD_TRANSPARENT_MODE:
		return "RD_TRANSPARENT_MODE"
	case CommonCommandWR_TX_ONLY_MODE:
		return "WR_TX_ONLY_MODE"
	case CommonCommandRD_TX_ONLY_MODE:
		return "RD_TX_ONLY_MODE"
	default:
		return "UNKNOWN"
	}
}

func (command CommonCommand) Valid() bool {
	switch command {
	case CommonCommandWR_SLEEP,
		CommonCommandWR_RESET,
		CommonCommandRD_VERSION,
		CommonCommandRD_SYS_LOG,
		CommonCommandRESET_SYS_LOG,
		CommonCommandWR_BIST,
		CommonCommandWR_IDBASE,
		CommonCommandRD_IDBASE,
		CommonCommandWR_REPEATER,
		CommonCommandRD_REPEATER,
		CommonCommandWR_FILTER_ADD,
		CommonCommandWR_FILTER_DEL,
		CommonCommandWR_FILTER_DEL_ALL,
		CommonCommandWR_FILTER_ENABLE,
		CommonCommandRD_FILTER,
		CommonCommandWR_WAIT_MATURITY,
		CommonCommandWR_SUBTEL,
		CommonCommandWR_MEM,
		CommonCommandRD_MEM,
		CommonCommandRD_MEM_ADDRESS,
		CommonCommandRD_SECURITY,
		CommonCommandWR_SECURITY,
		CommonCommandWR_LEARNMODE,
		CommonCommandRD_LEARNMODE,
		CommonCommandWR_SECUREDEVICE_ADD,
		CommonCommandWR_SECUREDEVICE_DEL,
		CommonCommandRD_SECUREDEVICE_BY_INDEX,
		CommonCommandWR_MODE,
		CommonCommandRD_NUMSECUREDEVICES,
		CommonCommandRD_SECUREDEVICE_BY_ID,
		CommonCommandWR_SECUREDEVICE_ADD_PSK,
		CommonCommandWR_SECUREDEVICE_SENDTEACHIN,
		CommonCommandWR_TEMPORARY_RLC_WINDOW,
		CommonCommandRD_SECUREDEVICE_PSK,
		CommonCommandRD_DUTYCYCLE_LIMIT,
		CommonCommandSET_BAUDRATE,
		CommonCommandGET_FREQUENCY_INFO,
		CommonCommandGET_STEPCODE,
		CommonCommandWR_REMAN_CODE,
		CommonCommandWR_STARTUP_DELAY,
		CommonCommandWR_REMAN_REPEATING,
		CommonCommandRD_REMAN_REPEATING,
		CommonCommandSET_NOISETHRESHOLD,
		CommonCommandGET_NOISETHRESHOLD,
		CommonCommandWR_RLC_SAVE_PERIOD,
		CommonCommandWR_RLC_LEGACY_MODE,
		CommonCommandWR_SECUREDEVICEV2_ADD,
		CommonCommandRD_SECUREDEVICEV2_BY_INDEX,
		CommonCommandWR_RSSITEST_MODE,
		CommonCommandRD_RSSITEST_MODE,
		CommonCommandWR_SECUREDEVICE_MAINTENANCEKEY,
		CommonCommandRD_SECUREDEVICE_MAINTENANCEKEY,
		CommonCommandWR_TRANSPARENT_MODE,
		CommonCommandRD_TRANSPARENT_MODE,
		CommonCommandWR_TX_ONLY_MODE,
		CommonCommandRD_TX_ONLY_MODE:
		return true
	default:
		return false
	}
}
