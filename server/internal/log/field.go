package mlog

import (
	"fmt"
)

type FieldType interface {
	any | string | bool |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func Field[T FieldType](msg string, value T) string {
	res := ""

	switch v := any(value).(type) {
	case error:
		res = fmt.Sprintf("%s, error: %s", msg, v.Error())
	case string:
		res = fmt.Sprintf("%s, %s", msg, v)
	case bool:
		res = fmt.Sprintf("%s, %t", msg, v)
	case float32, float64:
		res = fmt.Sprintf("%s, %.2f", msg, value)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		res = fmt.Sprintf("%s, %d", msg, v)
	default: // regard as struct
		res = fmt.Sprintf("%s. type: %T, value: %+v", msg, value, value)
	}

	return res
}
