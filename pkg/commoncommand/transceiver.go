package commoncommand

import (
	"errors"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

type RdDutyCycleLimit struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (cmd *RdDutyCycleLimit) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewRdDutyCycleLimit() (RdDutyCycleLimit, error) {
	return RdDutyCycleLimit{
		CommandCode: enums.CommonCommandRD_DUTYCYCLE_LIMIT,
	}, nil
}

type RdDutyCycleLimitResponse struct {
	AvailableDutyCycle                 uint8
	Slots                              uint8
	SlotPeriod                         uint16
	TimeLeftInCurrentSlot              uint16
	AvailableDutyCycleAfterCurrentSlot uint8
}

func ParseRdDutyCycleLimitResponseOK(response response.Packet) (RdDutyCycleLimitResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdDutyCycleLimitResponse{}, errors.New("invalid return code")
	}

	var result RdDutyCycleLimitResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdDutyCycleLimitResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type SetBaudrate struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	Baudrate    enums.TCMBaudrate   `enocean-esp3:"data"`
}

func (cmd *SetBaudrate) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewSetBaudrate(baudrate enums.TCMBaudrate) (SetBaudrate, error) {
	return SetBaudrate{
		CommandCode: enums.CommonCommandSET_BAUDRATE,
		Baudrate:    baudrate,
	}, nil
}

type GetFrequencyInfo struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (cmd *GetFrequencyInfo) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewGetFrequencyInfo() (GetFrequencyInfo, error) {
	return GetFrequencyInfo{
		CommandCode: enums.CommonCommandGET_FREQUENCY_INFO,
	}, nil
}

type GetFrequencyInfoResponse struct {
	Frequency enums.TCMFrequency
	Protocol  enums.TCMProtocol
}

func ParseGetFrequencyInfoResponseOK(response response.Packet) (GetFrequencyInfoResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return GetFrequencyInfoResponse{}, errors.New("invalid return code")
	}

	var result GetFrequencyInfoResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return GetFrequencyInfoResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type GetStepCode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (cmd *GetStepCode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewGetStepCode() (GetStepCode, error) {
	return GetStepCode{
		CommandCode: enums.CommonCommandGET_STEPCODE,
	}, nil
}

type GetStepCodeResponse struct {
	StepCode uint8
	Revision uint8
}

func ParseGetStepCodeResponseOK(response response.Packet) (GetStepCodeResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return GetStepCodeResponse{}, errors.New("invalid return code")
	}

	var result GetStepCodeResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return GetStepCodeResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type WrStartupDelay struct {
	CommandCode  enums.CommonCommand `enocean-esp3:"data"`
	StartupDelay uint8               `enocean-esp3:"data"` //Multiple of 10ms
}

func (cmd *WrStartupDelay) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewWrStartupDelay(startupDelay uint8) (WrStartupDelay, error) {
	return WrStartupDelay{
		CommandCode:  enums.CommonCommandWR_STARTUP_DELAY,
		StartupDelay: startupDelay,
	}, nil
}

type SetNoiseThreshold struct {
	CommandCode    enums.CommonCommand `enocean-esp3:"data"`
	NoiseThreshold uint8               `enocean-esp3:"data"`
}

func (cmd *SetNoiseThreshold) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewSetNoiseThreshold(noiseThreshold uint8) (SetNoiseThreshold, error) {
	return SetNoiseThreshold{
		CommandCode:    enums.CommonCommandSET_NOISETHRESHOLD,
		NoiseThreshold: noiseThreshold,
	}, nil
}

type GetNoiseThreshold struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (cmd *GetNoiseThreshold) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewGetNoiseThreshold() (GetNoiseThreshold, error) {
	return GetNoiseThreshold{
		CommandCode: enums.CommonCommandGET_NOISETHRESHOLD,
	}, nil
}

type GetNoiseThresholdResponse struct {
	RSSILevel uint8
}

func ParseGetNoiseThresholdResponseOK(response response.Packet) (GetNoiseThresholdResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return GetNoiseThresholdResponse{}, errors.New("invalid return code")
	}

	var result GetNoiseThresholdResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return GetNoiseThresholdResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type SetCRCMode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	CRCMode     enums.CRCMode       `enocean-esp3:"data"`
}

func (cmd *SetCRCMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewSetCRCMode(crcMode enums.CRCMode) (SetCRCMode, error) {
	return SetCRCMode{
		CommandCode: enums.CommonCommandSET_CRCMode,
		CRCMode:     crcMode,
	}, nil
}

type GetCRCMode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (cmd *GetCRCMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewGetCRCMode() (GetCRCMode, error) {
	return GetCRCMode{
		CommandCode: enums.CommonCommandGET_CRCMode,
	}, nil
}

type GetCRCModeResponse struct {
	CRCMode enums.CRCMode
}

func ParseGetCRCModeResponseOK(response response.Packet) (GetCRCModeResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return GetCRCModeResponse{}, errors.New("invalid return code")
	}

	var result GetCRCModeResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return GetCRCModeResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type WrRSSITestMode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	TestMode    enums.RSSITestMode  `enocean-esp3:"data"`
	Timeout     uint16              `enocean-esp3:"data"`
}

func (cmd *WrRSSITestMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewWrRSSITestMode(testMode enums.RSSITestMode, timeout uint16) (WrRSSITestMode, error) {
	return WrRSSITestMode{
		CommandCode: enums.CommonCommandWR_RSSITEST_MODE,
		TestMode:    testMode,
		Timeout:     timeout,
	}, nil
}

type RdRSSITestMode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (cmd *RdRSSITestMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewRdRSSITestMode() (RdRSSITestMode, error) {
	return RdRSSITestMode{
		CommandCode: enums.CommonCommandRD_RSSITEST_MODE,
	}, nil
}

type RdRSSITestModeResponse struct {
	TestMode enums.RSSITestMode `enocean-esp3:"data"`
}

func ParseRdRSSITestModeResponseOK(response response.Packet) (RdRSSITestModeResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdRSSITestModeResponse{}, errors.New("invalid return code")
	}

	var result RdRSSITestModeResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdRSSITestModeResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type WrTransparentMode struct {
	CommandCode     enums.CommonCommand   `enocean-esp3:"data"`
	TransparentMode enums.TransparentMode `enocean-esp3:"data"`
}

func (cmd *WrTransparentMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewWrTransparentMode(transparentMode enums.TransparentMode) (WrTransparentMode, error) {
	return WrTransparentMode{
		CommandCode:     enums.CommonCommandWR_TRANSPARENT_MODE,
		TransparentMode: transparentMode,
	}, nil
}

type RdTransparentMode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (cmd *RdTransparentMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewRdTransparentMode() (RdTransparentMode, error) {
	return RdTransparentMode{
		CommandCode: enums.CommonCommandRD_TRANSPARENT_MODE,
	}, nil
}

type RdTransparentModeResponse struct {
	TransparentMode enums.TransparentMode `enocean-esp3:"data"`
}

func ParseRdTransparentModeResponseOK(response response.Packet) (RdTransparentModeResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdTransparentModeResponse{}, errors.New("invalid return code")
	}

	var result RdTransparentModeResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdTransparentModeResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}

type WrTxOnlyMode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
	TxOnlyMode  enums.TxOnlyMode    `enocean-esp3:"data"`
}

func (cmd *WrTxOnlyMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewWrTxOnlyMode(txOnlyMode enums.TxOnlyMode) (WrTxOnlyMode, error) {
	return WrTxOnlyMode{
		CommandCode: enums.CommonCommandWR_TX_ONLY_MODE,
		TxOnlyMode:  txOnlyMode,
	}, nil
}

type RdTxOnlyMode struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (cmd *RdTxOnlyMode) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewRdTxOnlyMode() (RdTxOnlyMode, error) {
	return RdTxOnlyMode{
		CommandCode: enums.CommonCommandRD_TX_ONLY_MODE,
	}, nil
}

type RdTxOnlyModeResponse struct {
	TxOnlyMode enums.TxOnlyMode `enocean-esp3:"data"`
}

func ParseRdTxOnlyModeResponseOK(response response.Packet) (RdTxOnlyModeResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdTxOnlyModeResponse{}, errors.New("invalid return code")
	}

	var result RdTxOnlyModeResponse
	if err := serializer.BytesToStruct(response.Data, &result); err != nil {
		return RdTxOnlyModeResponse{}, errors.New("failed to deserialize response")
	}

	return result, nil
}
