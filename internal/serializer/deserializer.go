package serializer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

// CustomDeserializer is a function type for custom deserialization of specific
// types. It reads from the buffer and sets the value.
type CustomDeserializer func(buf *bytes.Reader, v reflect.Value, byteOrder binary.ByteOrder) error

// DeserializerConfig encapsulates configuration for deserialization.
type DeserializerConfig struct {
	// Deserializers is the per-call registry of custom deserializers keyed by
	// concrete reflect.Type.
	Deserializers map[reflect.Type]CustomDeserializer

	// ByteOrder controls how numeric values are decoded.
	// If nil, binary.BigEndian is used.
	ByteOrder binary.ByteOrder
}

// sanitize ensures we always have a non-nil map and byte order.
func (c DeserializerConfig) sanitize() DeserializerConfig {
	if c.Deserializers == nil {
		c.Deserializers = map[reflect.Type]CustomDeserializer{}
	}
	if c.ByteOrder == nil {
		c.ByteOrder = binary.BigEndian
	}
	return c
}

// merge combines multiple DeserializerConfig values into one.
// Later configs override earlier ones for the same fields.
func mergeDeserializerConfigs(configs []DeserializerConfig) DeserializerConfig {
	var merged DeserializerConfig

	for _, cfg := range configs {
		// Merge Deserializers maps
		if cfg.Deserializers != nil {
			if merged.Deserializers == nil {
				merged.Deserializers = make(map[reflect.Type]CustomDeserializer)
			}
			for k, v := range cfg.Deserializers {
				merged.Deserializers[k] = v
			}
		}

		// Use the last non-nil ByteOrder
		if cfg.ByteOrder != nil {
			merged.ByteOrder = cfg.ByteOrder
		}
	}

	return merged
}

// BytesToStruct deserializes a byte slice into a struct using reflection.
// The struct fields are read sequentially from the byte slice based on their types.
// All numeric types must have explicit sizes (uint8, uint16, uint32, uint64, not uint).
// The structPtr parameter must be a pointer to a struct.
// If no config is provided, defaults are used (BigEndian byte order, no custom deserializers).
// If multiple configs are provided, they are merged with later configs overriding earlier ones.
func BytesToStruct(data []byte, structPtr any, cfg ...DeserializerConfig) error {
	if structPtr == nil {
		return fmt.Errorf("structPtr cannot be nil")
	}

	v := reflect.ValueOf(structPtr)
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("structPtr must be a pointer to a struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("structPtr must point to a struct")
	}

	config := mergeDeserializerConfigs(cfg).sanitize()

	reader := bytes.NewReader(data)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Skip unexported fields
		if !fieldValue.CanSet() {
			continue
		}

		if err := deserializeValue(reader, fieldValue, config); err != nil {
			return fmt.Errorf("failed to deserialize field %s: %w", field.Name, err)
		}
	}

	return nil
}

// deserializeValue deserializes a value from the reader using the provided configuration.
func deserializeValue(reader *bytes.Reader, v reflect.Value, cfg DeserializerConfig) error {
	cfg = cfg.sanitize()

	// Handle pointer types
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			// Create a new value for the pointer
			elemType := v.Type().Elem()
			newVal := reflect.New(elemType).Elem()
			if err := deserializeValue(reader, newVal, cfg); err != nil {
				return err
			}
			v.Set(newVal.Addr())
			return nil
		}
		v = v.Elem()
	}

	// Check for custom deserializer in the provided configuration
	if customDeserializer, ok := cfg.Deserializers[v.Type()]; ok {
		return customDeserializer(reader, v, cfg.ByteOrder)
	}

	switch v.Kind() {
	case reflect.Struct:
		return deserializeStruct(reader, v, cfg)
	case reflect.Array:
		return deserializeArray(reader, v, cfg)
	case reflect.Slice:
		return deserializeSlice(reader, v, cfg)
	default:
		return binary.Read(reader, cfg.ByteOrder, v.Addr().Interface())
	}
}

// deserializeStruct deserializes a struct using the supplied configuration.
func deserializeStruct(reader *bytes.Reader, v reflect.Value, cfg DeserializerConfig) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		if err := deserializeValue(reader, fieldValue, cfg); err != nil {
			return fmt.Errorf("failed to deserialize nested field %s: %w", field.Name, err)
		}
	}

	return nil
}

// deserializeArray deserializes an array using the supplied configuration.
func deserializeArray(reader *bytes.Reader, v reflect.Value, cfg DeserializerConfig) error {
	cfg = cfg.sanitize()

	length := v.Len()
	elemType := v.Type().Elem()
	elemKind := elemType.Kind()

	// For byte arrays, read directly as bytes
	if elemKind == reflect.Uint8 {
		if reader.Len() < length {
			return fmt.Errorf("not enough bytes to read array of length %d (remaining: %d)", length, reader.Len())
		}
		bytes := make([]byte, length)
		if _, err := reader.Read(bytes); err != nil {
			return err
		}
		for i := 0; i < length; i++ {
			v.Index(i).SetUint(uint64(bytes[i]))
		}
		return nil
	}

	// For non-byte arrays, deserialize each element
	for i := 0; i < length; i++ {
		if err := deserializeValue(reader, v.Index(i), cfg); err != nil {
			return err
		}
	}
	return nil
}

// deserializeSlice deserializes a slice using the supplied configuration.
// For byte slices, reads all remaining bytes.
// For other slices, this requires knowing the length - which is a limitation.
// In practice, slices should have a length prefix or be the last field.
func deserializeSlice(reader *bytes.Reader, v reflect.Value, cfg DeserializerConfig) error {
	cfg = cfg.sanitize()

	elemType := v.Type().Elem()
	elemKind := elemType.Kind()

	// For byte slices, read all remaining bytes
	if elemKind == reflect.Uint8 {
		remaining := reader.Len()
		if remaining > 0 {
			bytes := make([]byte, remaining)
			n, err := reader.Read(bytes)
			if err != nil {
				// If we read some bytes before hitting EOF, that's okay
				if n == 0 {
					return err
				}
			}
			v.SetBytes(bytes[:n])
		} else {
			v.SetBytes(nil)
		}
		return nil
	}

	// Non-byte slices must be the last field; consume complete elements until
	// the input is exhausted and reject malformed trailing data.
	slice := reflect.MakeSlice(v.Type(), 0, 0)
	for reader.Len() > 0 {
		before := reader.Len()
		elem := reflect.New(elemType).Elem()
		if err := deserializeValue(reader, elem, cfg); err != nil {
			return err
		}
		if reader.Len() == before {
			return fmt.Errorf("deserializer made no progress for %s", elemType)
		}
		slice = reflect.Append(slice, elem)
	}
	v.Set(slice)
	return nil
}
