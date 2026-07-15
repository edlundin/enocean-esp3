package serializer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"

	"github.com/edlundin/enocean-esp3/pkg/enums"
	"github.com/edlundin/enocean-esp3/pkg/esp3"
)

// CustomSerializer is a function type for custom serialization of specific
// types in the pure serializer. It is identical in spirit to CustomSerializer
// but does not rely on any global registry.
type CustomSerializer func(buf *bytes.Buffer, v reflect.Value, byteOrder binary.ByteOrder) error

// SerializerConfig encapsulates all configuration required to perform
// serialization in a pure/functional style – there is no global mutable state.
type SerializerConfig struct {
	// Serializers is the per-call registry of custom serializers keyed by
	// concrete reflect.Type.
	Serializers map[reflect.Type]CustomSerializer

	// ByteOrder controls how numeric values are encoded.
	// If nil, binary.BigEndian is used.
	ByteOrder binary.ByteOrder
}

// sanitize ensures we always have a non-nil map and byte order.
func (c SerializerConfig) sanitize() SerializerConfig {
	if c.Serializers == nil {
		c.Serializers = map[reflect.Type]CustomSerializer{}
	}
	if c.ByteOrder == nil {
		c.ByteOrder = binary.BigEndian
	}
	return c
}

// merge combines multiple SerializerConfig values into one.
// Later configs override earlier ones for the same fields.
// Serializers maps are merged (later entries override earlier ones for the same type).
// ByteOrder uses the last non-nil value, or defaults to BigEndian if all are nil.
func mergeConfigs(configs []SerializerConfig) SerializerConfig {
	var merged SerializerConfig

	for _, cfg := range configs {
		// Merge Serializers maps
		if cfg.Serializers != nil {
			if merged.Serializers == nil {
				merged.Serializers = make(map[reflect.Type]CustomSerializer)
			}
			for k, v := range cfg.Serializers {
				merged.Serializers[k] = v
			}
		}

		// Use the last non-nil ByteOrder
		if cfg.ByteOrder != nil {
			merged.ByteOrder = cfg.ByteOrder
		}
	}

	return merged
}

// CommandToTelegram is a pure/functional variant of CommandToTelegram.
// All behavior is driven by the supplied SerializerConfig; there is no
// global registry or mutable package-level state.
// If no config is provided, defaults are used (BigEndian byte order, no custom serializers).
// If multiple configs are provided, they are merged with later configs overriding earlier ones.
func CommandToTelegram(cmd any, cfg ...SerializerConfig) (esp3.Telegram, error) {
	config := mergeConfigs(cfg).sanitize()

	v := reflect.ValueOf(cmd)

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return esp3.Telegram{}, fmt.Errorf("command must be a struct")
	}

	bufData := bytes.NewBuffer(nil)
	bufOptData := bytes.NewBuffer(nil)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !fieldValue.CanInterface() {
			continue
		}

		target := ""
		skipIfNone := false

		if tag := field.Tag.Get("enocean-esp3"); tag != "" {
			parts := strings.Split(tag, ",")
			target = parts[0]
			skipIfNone = len(parts) > 1 && parts[1] == "skipif:none"
		} else {
			continue
		}

		if skipIfNone && fieldValue.IsZero() {
			continue
		}

		var buf *bytes.Buffer
		switch target {
		case "data":
			buf = bufData
		case "optdata":
			buf = bufOptData
		default:
			return esp3.Telegram{}, fmt.Errorf("invalid target: %s", target)
		}

		if err := serializeValue(buf, fieldValue, config); err != nil {
			return esp3.Telegram{}, fmt.Errorf("failed to serialize field %s: %w", field.Name, err)
		}
	}

	optData := bufOptData.Bytes()
	if len(optData) == 0 {
		optData = nil
	}

	return esp3.NewTelegramFromData(enums.PacketTypeCOMMON_COMMAND, bufData.Bytes(), optData), nil
}

// serializeValue serializes a value using the provided configuration.
func serializeValue(buf *bytes.Buffer, v reflect.Value, cfg SerializerConfig) error {
	cfg = cfg.sanitize()

	// Handle pointer and interface types
	if v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			var zeroBytes []byte

			if v.Kind() == reflect.Pointer {
				zeroBytes = make([]byte, v.Type().Elem().Size())
			} else {
				return fmt.Errorf("nil interface: size unknown")
			}

			_, err := buf.Write(zeroBytes)

			return err
		}

		v = v.Elem()
	}

	// Check for custom serializer in the provided configuration
	if customSerializer, ok := cfg.Serializers[v.Type()]; ok {
		return customSerializer(buf, v, cfg.ByteOrder)
	}

	switch v.Kind() {
	case reflect.Struct:
		return serializeStruct(buf, v, cfg)
	case reflect.Array, reflect.Slice:
		return serializeSequence(buf, v, cfg)
	case reflect.String:
		_, err := buf.WriteString(v.String())
		return err
	default:
		return binary.Write(buf, cfg.ByteOrder, v.Interface())
	}
}

// serializeStruct serializes a struct using the supplied configuration.
func serializeStruct(buf *bytes.Buffer, v reflect.Value, cfg SerializerConfig) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		// Skip fields with enocean-esp3 tags – handled at the top level.
		if field.Tag.Get("enocean-esp3") != "" {
			continue
		}

		if err := serializeValue(buf, fieldValue, cfg); err != nil {
			return fmt.Errorf("failed to serialize nested field %s: %w", field.Name, err)
		}
	}

	return nil
}

// serializeSequence serializes arrays and slices using the supplied configuration.
func serializeSequence(buf *bytes.Buffer, v reflect.Value, cfg SerializerConfig) error {
	cfg = cfg.sanitize()

	kind := v.Kind()
	length := v.Len()

	if kind == reflect.Slice && v.IsNil() {
		return nil
	}

	elemKind := v.Type().Elem().Kind()

	if elemKind == reflect.Uint8 {
		if kind == reflect.Slice {
			_, err := buf.Write(v.Bytes())
			return err
		}

		data := make([]byte, length)
		for i := range data {
			data[i] = byte(v.Index(i).Uint())
		}
		_, err := buf.Write(data)

		return err
	}

	for i := 0; i < length; i++ {
		if err := serializeValue(buf, v.Index(i), cfg); err != nil {
			return err
		}
	}

	return nil
}
