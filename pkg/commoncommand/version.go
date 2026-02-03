package commoncommand

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"

	"github.com/edlundin/enocean-esp3/internal/serializer"
	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
	"github.com/edlundin/enocean-esp3/pkg/response"
)

type RdVersion struct {
	CommandCode enums.CommonCommand `enocean-esp3:"data"`
}

func (cmd *RdVersion) Serialize() (esp3.Telegram, error) {
	return serializer.CommandToTelegram(cmd)
}

func NewRdVersion() (RdVersion, error) {
	return RdVersion{
		CommandCode: enums.CommonCommandRD_VERSION,
	}, nil
}

type RdVersionResponse struct {
	AppVersion  [4]byte
	ApiVersion  [4]byte
	ChipID      uint32
	ChipVersion uint32
	Description string
}

func ParseRdVersionResponseOK(response response.Packet) (RdVersionResponse, error) {
	if response.Code != enums.ReturnCodeSUCCESS {
		return RdVersionResponse{}, errors.New("invalid return code")
	}

	// Use custom deserializer to handle string
	stringDeserializer := func(buf *bytes.Reader, v reflect.Value, _ binary.ByteOrder) error {
		rest, _ := io.ReadAll(buf) // io.ReadAll on bytes.Reader never fails
		v.SetString(string(rest))
		return nil
	}

	cfg := serializer.DeserializerConfig{
		Deserializers: map[reflect.Type]serializer.CustomDeserializer{
			reflect.TypeOf(""): stringDeserializer,
		},
	}

	var version RdVersionResponse
	if err := serializer.BytesToStruct(response.Data, &version, cfg); err != nil {
		return RdVersionResponse{}, err
	}

	return version, nil
}
