package commoncommand

import (
	"errors"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

// RdSysLog is a command to read the system logs
type RdSysLog struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

// Serialize encodes RdSysLog into its wire representation.
func (cmd *RdSysLog) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewRdSysLog constructs RdSysLog.
func NewRdSysLog() (RdSysLog, error) {
	return RdSysLog{
		CommandCode: enums.CommonCommandRD_SYS_LOG,
	}, nil
}

type RdSysLogResponse struct {
	ApiLogEntries []byte
	AppLogEntries []byte
}

// ParseRdSysLogResponseOK parses RdSysLogResponseOK.
func ParseRdSysLogResponseOK(response response.Packet) (RdSysLogResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdSysLogResponse{}, errors.New("invalid return code")
	}

	return RdSysLogResponse{
		ApiLogEntries: response.Data,
		AppLogEntries: response.OptData,
	}, nil
}

// ResetSysLog is a command to reset the system logs
type ResetSysLog struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

// Serialize encodes ResetSysLog into its wire representation.
func (cmd *ResetSysLog) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

// NewResetSysLog constructs ResetSysLog.
func NewResetSysLog() (ResetSysLog, error) {
	return ResetSysLog{
		CommandCode: enums.CommonCommandRESET_SYS_LOG,
	}, nil
}
