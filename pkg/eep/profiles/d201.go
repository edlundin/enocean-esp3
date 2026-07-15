package profiles

import "fmt"

// D201SetOutput encodes the D2-01-00 SET_OUTPUT command.
func D201SetOutput(channel, value uint8) ([]byte, error) {
	if channel == 0 || channel > 31 {
		return nil, fmt.Errorf("D2-01 channel must be 1..31")
	}
	if value > 100 {
		return nil, fmt.Errorf("D2-01 output must be 0..100")
	}
	return []byte{1, channel, value}, nil
}

type D201Status struct {
	Channel uint8
	Output  uint8
}

// ParseD201Status decodes the D2-01-00 STATUS_RESPONSE output value.
func ParseD201Status(data []byte) (D201Status, error) {
	if len(data) < 3 {
		return D201Status{}, fmt.Errorf("D2-01 status response must be at least 3 bytes")
	}
	if data[0]&0x0f != 4 {
		return D201Status{}, fmt.Errorf("unsupported D2-01 command %d", data[0]&0x0f)
	}
	channel, output := data[1]&0x1f, data[2]&0x7f
	if channel == 0 || output > 100 {
		return D201Status{}, fmt.Errorf("unsupported D2-01 status response")
	}
	return D201Status{Channel: channel, Output: output}, nil
}
