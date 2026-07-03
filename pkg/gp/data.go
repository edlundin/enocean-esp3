package gp

import "errors"

type SelectedValue struct {
	Index int
	Value uint64
}

func EncodeCompleteData(channels []Channel, values []uint64) ([]byte, error) {
	ops, err := operationalChannels(channels)
	if err != nil {
		return nil, err
	}
	if len(values) != len(ops) {
		return nil, errors.New("GP complete data value count does not match channels")
	}
	bits := 0
	widths := make([]int, len(ops))
	for i, ch := range ops {
		w, err := ch.ValueBits()
		if err != nil {
			return nil, err
		}
		widths[i] = w
		bits += w
	}
	out := bytesForBits(bits)
	if len(out) > MaxMessageLength {
		return nil, errors.New("GP complete data exceeds max message length")
	}
	off := 0
	for i, v := range values {
		if err := writeUnsigned(out, off, widths[i], v); err != nil {
			return nil, err
		}
		off += widths[i]
	}
	return out, nil
}

func DecodeCompleteData(channels []Channel, data []byte) ([]uint64, error) {
	ops, err := operationalChannels(channels)
	if err != nil {
		return nil, err
	}
	if len(data) > MaxMessageLength {
		return nil, errors.New("GP complete data exceeds max message length")
	}
	out := make([]uint64, len(ops))
	off := 0
	for i, ch := range ops {
		w, err := ch.ValueBits()
		if err != nil {
			return nil, err
		}
		out[i], err = readUnsigned(data, off, w)
		if err != nil {
			return nil, err
		}
		off += w
	}
	return out, nil
}

func EncodeSelectiveData(channels []Channel, values []SelectedValue) ([]byte, error) {
	ops, err := operationalChannels(channels)
	if err != nil {
		return nil, err
	}
	if len(values) > 15 {
		return nil, errors.New("GP selective data supports at most 15 channels")
	}
	bits := 4
	widths := make([]int, len(values))
	for i, sv := range values {
		if sv.Index < 0 || sv.Index >= len(ops) || sv.Index > 63 {
			return nil, errors.New("GP selective data channel index out of range")
		}
		w, err := ops[sv.Index].ValueBits()
		if err != nil {
			return nil, err
		}
		widths[i] = w
		bits += 6 + w
	}
	out := bytesForBits(bits)
	if len(out) > MaxMessageLength {
		return nil, errors.New("GP selective data exceeds max message length")
	}
	_ = writeUnsigned(out, 0, 4, uint64(len(values)))
	off := 4
	for i, sv := range values {
		_ = writeUnsigned(out, off, 6, uint64(sv.Index))
		off += 6
		_ = writeUnsigned(out, off, widths[i], sv.Value)
		off += widths[i]
	}
	return out, nil
}

func DecodeSelectiveData(channels []Channel, data []byte) ([]SelectedValue, error) {
	ops, err := operationalChannels(channels)
	if err != nil {
		return nil, err
	}
	count, err := readUnsigned(data, 0, 4)
	if err != nil {
		return nil, err
	}
	out := make([]SelectedValue, int(count))
	off := 4
	for i := range out {
		idx, err := readUnsigned(data, off, 6)
		if err != nil {
			return nil, err
		}
		off += 6
		if int(idx) >= len(ops) {
			return nil, errors.New("GP selective data channel index out of range")
		}
		w, err := ops[idx].ValueBits()
		if err != nil {
			return nil, err
		}
		v, err := readUnsigned(data, off, w)
		if err != nil {
			return nil, err
		}
		off += w
		out[i] = SelectedValue{Index: int(idx), Value: v}
	}
	return out, nil
}

func operationalChannels(channels []Channel) ([]Channel, error) {
	out := make([]Channel, 0, len(channels))
	for _, ch := range channels {
		if ch.Type == ChannelTeachInInformation {
			continue
		}
		if _, err := ch.ValueBits(); err != nil {
			return nil, err
		}
		out = append(out, ch)
	}
	return out, nil
}
