package generic_value

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ToGenericValue converts various types to a generic interface{} value
func ToGenericValue(value interface{}) interface{} {
	return value
}

// ToString converts a value to string
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	
	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToInt converts a value to int
func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(strings.TrimSpace(v))
	default:
		return 0, fmt.Errorf("cannot convert %v (type %T) to int", value, value)
	}
}

// GetReflectValue returns the reflect.Value of an interface
func GetReflectValue(value interface{}) reflect.Value {
	return reflect.ValueOf(value)
}

// SedulousTypeConverter provides methods for converting between types
type SedulousTypeConverter struct{}

// ConvertToString converts a value to string
func (s *SedulousTypeConverter) ConvertToString(value interface{}) string {
	return ToString(value)
}