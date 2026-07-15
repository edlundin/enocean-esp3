package commoncommand

import (
	"errors"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/deviceid"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

// NOTE: For newer devices (e.g. TCM515), use WrSecureDeviceV2Add instead
type WrSecureDeviceAdd struct {
	CommandCode         enums.CommonCommand         `enocean-esp3:"data"`
	SecurityLevelFormat uint8                       `enocean-esp3:"data"`
	DeviceID            deviceid.DeviceID           `enocean-esp3:"data"`
	SecurityKey         [16]byte                    `enocean-esp3:"data"`
	RollingCode         [3]byte                     `enocean-esp3:"data"`
	Direction           enums.SecureDeviceDirection `enocean-esp3:"optdata"`
	PTMModule           uint8                       `enocean-esp3:"optdata"`
	TeachInInfo         uint8                       `enocean-esp3:"optdata"`
}

// Serialize encodes WrSecureDeviceAdd into its wire representation.
func (cmd *WrSecureDeviceAdd) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrSecureDeviceAdd constructs WrSecureDeviceAdd.
func NewWrSecureDeviceAdd(securityLevelFormat uint8, deviceID deviceid.DeviceID, securityKey [16]byte, rollingCode [3]byte, direction enums.SecureDeviceDirection, ptmModule uint8, teachInfo uint8) (WrSecureDeviceAdd, error) {
	if teachInfo > 0x0f {
		return WrSecureDeviceAdd{}, errors.New("teach info out of range: only half a byte is allowed, use NewWrSecureDeviceV2Add 1-byte teach-in info")
	}

	return WrSecureDeviceAdd{
		CommandCode:         enums.CommonCommandWR_SECUREDEVICE_ADD,
		SecurityLevelFormat: securityLevelFormat,
		DeviceID:            deviceID,
		SecurityKey:         securityKey,
		RollingCode:         rollingCode,
		Direction:           direction,
		PTMModule:           ptmModule,
		TeachInInfo:         teachInfo,
	}, nil
}

type WrSecureDeviceDel struct {
	CommandCode enums.CommonCommand         `enocean-esp3:"data"`
	DeviceID    deviceid.DeviceID           `enocean-esp3:"data"`
	Direction   enums.SecureDeviceDirection `enocean-esp3:"optdata"`
}

// Serialize encodes WrSecureDeviceDel into its wire representation.
func (cmd *WrSecureDeviceDel) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrSecureDeviceDel constructs WrSecureDeviceDel.
func NewWrSecureDeviceDel(deviceID deviceid.DeviceID, direction enums.SecureDeviceDirection) (WrSecureDeviceDel, error) {
	return WrSecureDeviceDel{
		CommandCode: enums.CommonCommandWR_SECUREDEVICE_DEL,
		DeviceID:    deviceID,
		Direction:   direction,
	}, nil
}

type RdSecureDeviceByIndex struct {
	CommandCode enums.CommonCommand         `enocean-esp3:"data"`
	Index       uint8                       `enocean-esp3:"data"`
	Direction   enums.SecureDeviceDirection `enocean-esp3:"optdata"`
}

// Serialize encodes RdSecureDeviceByIndex into its wire representation.
func (cmd *RdSecureDeviceByIndex) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewRdSecureDeviceByIndex constructs RdSecureDeviceByIndex.
func NewRdSecureDeviceByIndex(index uint8, direction enums.SecureDeviceDirection) (RdSecureDeviceByIndex, error) {
	if index > 0xfe {
		return RdSecureDeviceByIndex{}, errors.New("index must be between 0 and 254")
	}

	return RdSecureDeviceByIndex{
		CommandCode: enums.CommonCommandRD_SECUREDEVICE_BY_INDEX,
		Index:       index,
		Direction:   direction,
	}, nil
}

type RdSecureDeviceByIndexResponse struct {
	SecurityLevelFormat uint8
	DeviceID            deviceid.DeviceID
	PrivateKey          [16]byte
	RollingCode         [3]byte
	PSK                 [16]byte
	TeachInInfo         uint8
}

// ParseRdSecureDeviceByIndexResponseOK parses RdSecureDeviceByIndexResponseOK.
func ParseRdSecureDeviceByIndexResponseOK(response response.Packet) (RdSecureDeviceByIndexResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdSecureDeviceByIndexResponse{}, errors.New("invalid return code")
	}

	mergedData := make([]byte, 0, len(response.Data)+len(response.OptData))
	mergedData = append(mergedData, response.Data...)
	mergedData = append(mergedData, response.OptData...)

	var result RdSecureDeviceByIndexResponse
	if err := serializer.BytesToStruct(mergedData, &result); err != nil {
		return RdSecureDeviceByIndexResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type RdNumSecureDevices struct {
	CommandCode enums.CommonCommand         `enocean-esp3:"data"`
	Direction   enums.SecureDeviceDirection `enocean-esp3:"optdata"`
}

// Serialize encodes RdNumSecureDevices into its wire representation.
func (cmd *RdNumSecureDevices) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewRdNumSecureDevices constructs RdNumSecureDevices.
func NewRdNumSecureDevices(direction enums.SecureDeviceDirection) (RdNumSecureDevices, error) {
	return RdNumSecureDevices{
		CommandCode: enums.CommonCommandRD_NUMSECUREDEVICES,
		Direction:   direction,
	}, nil
}

type RdNumSecureDevicesResponse struct {
	NumSecureDevices uint8
}

// ParseRdNumSecureDevicesResponseOK parses RdNumSecureDevicesResponseOK.
func ParseRdNumSecureDevicesResponseOK(response response.Packet) (RdNumSecureDevicesResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdNumSecureDevicesResponse{}, errors.New("invalid return code")
	}

	var result RdNumSecureDevicesResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdNumSecureDevicesResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type RdSecureDeviceByID struct {
	CommandCode enums.CommonCommand         `enocean-esp3:"data"`
	DeviceID    deviceid.DeviceID           `enocean-esp3:"data"`
	Direction   enums.SecureDeviceDirection `enocean-esp3:"optdata"`
}

// Serialize encodes RdSecureDeviceByID into its wire representation.
func (cmd *RdSecureDeviceByID) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewRdSecureDeviceByID constructs RdSecureDeviceByID.
func NewRdSecureDeviceByID(deviceID deviceid.DeviceID, direction enums.SecureDeviceDirection) (RdSecureDeviceByID, error) {
	return RdSecureDeviceByID{
		CommandCode: enums.CommonCommandRD_SECUREDEVICE_BY_ID,
		DeviceID:    deviceID,
		Direction:   direction,
	}, nil
}

type RdSecureDeviceByIDResponse struct {
	SecurityLevelFormat uint8
	Index               uint8
}

// ParseRdSecureDeviceByIDResponseOK parses RdSecureDeviceByIDResponseOK.
func ParseRdSecureDeviceByIDResponseOK(response response.Packet) (RdSecureDeviceByIDResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdSecureDeviceByIDResponse{}, errors.New("invalid return code")
	}

	var result RdSecureDeviceByIDResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdSecureDeviceByIDResponse{}, errors.New("failed to deserialize response")
	}

	if result.Index > 0xfe {
		return RdSecureDeviceByIDResponse{}, errors.New("index out of range")
	}

	return result, nil
}

type WrSecureDeviceAddPSK struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	DeviceID    deviceid.DeviceID   `enocean-esp3:"data"`
	PSK         [16]byte            `enocean-esp3:"data"`
}

// Serialize encodes WrSecureDeviceAddPSK into its wire representation.
func (cmd *WrSecureDeviceAddPSK) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrSecureDeviceAddPSK constructs WrSecureDeviceAddPSK.
func NewWrSecureDeviceAddPSK(deviceID deviceid.DeviceID, psk [16]byte) (WrSecureDeviceAddPSK, error) {
	return WrSecureDeviceAddPSK{
		CommandCode: enums.CommonCommandWR_SECUREDEVICE_ADD_PSK,
		DeviceID:    deviceID,
		PSK:         psk,
	}, nil
}

type WrSecureDeviceSendTeachIn struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	DeviceID    deviceid.DeviceID   `enocean-esp3:"data"`
	TeachInInfo uint8               `enocean-esp3:"optdata"`
}

// Serialize encodes WrSecureDeviceSendTeachIn into its wire representation.
func (cmd *WrSecureDeviceSendTeachIn) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrSecureDeviceSendTeachIn constructs WrSecureDeviceSendTeachIn.
func NewWrSecureDeviceSendTeachIn(deviceID deviceid.DeviceID, teachInInfo uint8) (WrSecureDeviceSendTeachIn, error) {
	return WrSecureDeviceSendTeachIn{
		CommandCode: enums.CommonCommandWR_SECUREDEVICE_SENDTEACHIN,
		DeviceID:    deviceID,
		TeachInInfo: teachInInfo,
	}, nil
}

type WrTemporaryRLCWindow struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	Enable      bool                `enocean-esp3:"data"`
	RLCWindow   uint32              `enocean-esp3:"data"`
}

// Serialize encodes WrTemporaryRLCWindow into its wire representation.
func (cmd *WrTemporaryRLCWindow) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrTemporaryRLCWindow constructs WrTemporaryRLCWindow.
func NewWrTemporaryRLCWindow(enable bool, rlcWindow uint32) (WrTemporaryRLCWindow, error) {
	return WrTemporaryRLCWindow{
		CommandCode: enums.CommonCommandWR_TEMPORARY_RLC_WINDOW,
		Enable:      enable,
		RLCWindow:   rlcWindow,
	}, nil
}

type RdSecureDevicePSK struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	DeviceID    deviceid.DeviceID   `enocean-esp3:"data"` // Use 0x00000000 for current device, other IDs will return the PSK for the device with the given ID
}

// Serialize encodes RdSecureDevicePSK into its wire representation.
func (cmd *RdSecureDevicePSK) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewRdSecureDevicePSK constructs RdSecureDevicePSK.
func NewRdSecureDevicePSK(deviceID deviceid.DeviceID, direction enums.SecureDeviceDirection) (RdSecureDevicePSK, error) {
	return RdSecureDevicePSK{
		CommandCode: enums.CommonCommandRD_SECUREDEVICE_PSK,
		DeviceID:    deviceID,
	}, nil
}

type RdSecureDevicePSKResponse struct {
	PSK [16]byte
}

// ParseRdSecureDevicePSKResponseOK parses RdSecureDevicePSKResponseOK.
func ParseRdSecureDevicePSKResponseOK(response response.Packet) (RdSecureDevicePSKResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdSecureDevicePSKResponse{}, errors.New("invalid return code")
	}

	var result RdSecureDevicePSKResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdSecureDevicePSKResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type WrRLCSavePeriod struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	SavePeriod  uint8               `enocean-esp3:"data"`
}

// Serialize encodes WrRLCSavePeriod into its wire representation.
func (cmd *WrRLCSavePeriod) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrRLCSavePeriod constructs WrRLCSavePeriod.
func NewWrRLCSavePeriod(savePeriod uint8) (WrRLCSavePeriod, error) {
	return WrRLCSavePeriod{
		CommandCode: enums.CommonCommandWR_RLC_SAVE_PERIOD,
		SavePeriod:  savePeriod,
	}, nil
}

type WrRLCLegacyMode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	RLCMode     enums.RLCMode       `enocean-esp3:"data"`
}

// Serialize encodes WrRLCLegacyMode into its wire representation.
func (cmd *WrRLCLegacyMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewWrRLCLegacyMode constructs WrRLCLegacyMode.
func NewWrRLCLegacyMode(rlcMode enums.RLCMode) (WrRLCLegacyMode, error) {
	if !rlcMode.Valid() {
		return WrRLCLegacyMode{}, errors.New("invalid RLC mode")
	}

	return WrRLCLegacyMode{
		CommandCode: enums.CommonCommandWR_RLC_LEGACY_MODE,
		RLCMode:     rlcMode,
	}, nil
}

type WrSecureDeviceV2Add struct {
	CommandCode         enums.CommonCommand         `enocean-esp3:"data"`
	SecurityLevelFormat uint8                       `enocean-esp3:"data"`
	DeviceID            deviceid.DeviceID           `enocean-esp3:"data"`
	PrivateKey          [16]byte                    `enocean-esp3:"data"`
	RollingCode         uint32                      `enocean-esp3:"data"`
	Direction           enums.SecureDeviceDirection `enocean-esp3:"optdata"`
}

// Serialize encodes WrSecureDeviceV2Add into its wire representation.
func (cmd *WrSecureDeviceV2Add) Serialize() (esp3.Telegram, error) {
	if cmd.Direction != enums.SecureDeviceDirectionINBOUND_TABLE &&
		cmd.Direction != enums.SecureDeviceDirectionOUTBOUND_TABLE &&
		cmd.Direction != enums.SecureDeviceDirectionOUTBOUND_BROADCAST_TABLE {
		return esp3.Telegram{}, errors.New("direction must be INBOUND_TABLE, OUTBOUND_TABLE or OUTBOUND_BROADCAST_TABLE")
	}

	return serializer.CommandToTelegram(cmd)
}

// NewWrSecureDeviceV2Add constructs WrSecureDeviceV2Add.
func NewWrSecureDeviceV2Add(securityLevelFormat uint8, deviceID deviceid.DeviceID, privateKey [16]byte, rollingCode uint32, direction enums.SecureDeviceDirection) (WrSecureDeviceV2Add, error) {
	return WrSecureDeviceV2Add{
		CommandCode:         enums.CommonCommandWR_SECUREDEVICEV2_ADD,
		SecurityLevelFormat: securityLevelFormat,
		DeviceID:            deviceID,
		PrivateKey:          privateKey,
		RollingCode:         rollingCode,
		Direction:           direction,
	}, nil
}

type RdSecureDeviceV2ByIndex struct {
	CommandCode enums.CommonCommand         `enocean-esp3:"data"`
	Index       uint8                       `enocean-esp3:"data"`
	Direction   enums.SecureDeviceDirection `enocean-esp3:"optdata"`
}

// Serialize encodes RdSecureDeviceV2ByIndex into its wire representation.
func (cmd *RdSecureDeviceV2ByIndex) Serialize() (esp3.Telegram, error) {
	if cmd.Index > 0xfe {
		return esp3.Telegram{}, errors.New("index must be between 0 and 254")
	}

	if cmd.Direction != enums.SecureDeviceDirectionINBOUND_TABLE && cmd.Direction != enums.SecureDeviceDirectionOUTBOUND_TABLE && cmd.Direction != enums.SecureDeviceDirectionOUTBOUND_BROADCAST_TABLE {
		return esp3.Telegram{}, errors.New("direction must be INBOUND_TABLE, OUTBOUND_TABLE or OUTBOUND_BROADCAST_TABLE")
	}

	return serializer.CommandToTelegram(cmd)
}

// NewRdSecureDeviceV2ByIndex constructs RdSecureDeviceV2ByIndex.
func NewRdSecureDeviceV2ByIndex(index uint8, direction enums.SecureDeviceDirection) (RdSecureDeviceV2ByIndex, error) {
	if index > 0xfe {
		return RdSecureDeviceV2ByIndex{}, errors.New("index must be between 0 and 254")
	}
	return RdSecureDeviceV2ByIndex{
		CommandCode: enums.CommonCommandRD_SECUREDEVICEV2_BY_INDEX,
		Index:       index,
		Direction:   direction,
	}, nil
}

type RdSecureDeviceV2ByIndexResponse struct {
	SecurityLevelFormat uint8
	DeviceID            deviceid.DeviceID
	PrivateKey          [16]byte
	RollingCode         uint32
	TeachInInfo         uint8
	PSK                 [16]byte
}

// ParseRdSecureDeviceV2ByIndexResponseOK parses RdSecureDeviceV2ByIndexResponseOK.
func ParseRdSecureDeviceV2ByIndexResponseOK(response response.Packet) (RdSecureDeviceV2ByIndexResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdSecureDeviceV2ByIndexResponse{}, errors.New("invalid return code")
	}

	mergedData := make([]byte, 0, len(response.Data)+len(response.OptData))
	mergedData = append(mergedData, response.Data...)
	mergedData = append(mergedData, response.OptData...)

	var result RdSecureDeviceV2ByIndexResponse
	if err := serializer.BytesToStruct(mergedData, &result); err != nil {
		return RdSecureDeviceV2ByIndexResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type WrSecureDeviceRemanKey struct {
	CommandCode    enums.CommonCommand `enocean-esp3:"data"`
	DeviceID       deviceid.DeviceID   `enocean-esp3:"data"`
	RemanKey       [16]byte            `enocean-esp3:"data"`
	RemanKeyNumber uint8               `enocean-esp3:"data"`
}

// Serialize encodes WrSecureDeviceRemanKey into its wire representation.
func (cmd *WrSecureDeviceRemanKey) Serialize() (esp3.Telegram, error) {
	if cmd.RemanKeyNumber < 1 || cmd.RemanKeyNumber > 0x0f {
		return esp3.Telegram{}, errors.New("reman key number must be between 1 and 15")
	}

	return serializer.CommandToTelegram(cmd)
}

// NewWrSecureDeviceRemainCode constructs a secure-device ReMan key command.
func NewWrSecureDeviceRemainCode(deviceID deviceid.DeviceID, remanKey [16]byte, remanKeyNumber uint8) (WrSecureDeviceRemanKey, error) {
	return WrSecureDeviceRemanKey{
		CommandCode:    enums.CommonCommandWR_SECUREDEVICE_REMAN_KEY,
		DeviceID:       deviceID,
		RemanKey:       remanKey,
		RemanKeyNumber: remanKeyNumber,
	}, nil
}

type RdSecureDeviceRemanKey struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	Index       uint8               `enocean-esp3:"data"`
}

// Serialize encodes RdSecureDeviceRemanKey into its wire representation.
func (cmd *RdSecureDeviceRemanKey) Serialize() (esp3.Telegram, error) {
	if cmd.Index < 1 || cmd.Index > 0x0f {
		return esp3.Telegram{}, errors.New("index must be between 1 and 15")
	}

	return serializer.CommandToTelegram(cmd)
}

// NewRdSecureDeviceRemanKey constructs RdSecureDeviceRemanKey.
func NewRdSecureDeviceRemanKey(index uint8) (RdSecureDeviceRemanKey, error) {
	return RdSecureDeviceRemanKey{
		CommandCode: enums.CommonCommandRD_SECUREDEVICE_REMAN_KEY,
		Index:       index,
	}, nil
}

type RdSecureDeviceRemanKeyResponse struct {
	Index               uint8
	DeviceID            deviceid.DeviceID
	PrivateKey          [16]byte
	KeyNumber           uint8
	InboundRollingCode  uint32
	OutboundRollingCode uint32
}

// ParseRdSecureDeviceRemanKeyResponseOK parses RdSecureDeviceRemanKeyResponseOK.
func ParseRdSecureDeviceRemanKeyResponseOK(response response.Packet) (RdSecureDeviceRemanKeyResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdSecureDeviceRemanKeyResponse{}, errors.New("invalid return code")
	}

	var result RdSecureDeviceRemanKeyResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdSecureDeviceRemanKeyResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}
