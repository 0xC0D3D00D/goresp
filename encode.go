package goresp

import (
	"errors"
	"strconv"
)

var (
	ErrUnsupportedType = errors.New("Unsupported Type")
	crlf               = []byte{'\r', '\n'}
)

func encodeNil() []byte {
	return []byte{'$', '-', '1', '\r', '\n'}
}

func encodeError(e error) []byte {
	err := e.Error()
	encoded := append([]byte{'-'}, err...)
	encoded = append(encoded, crlf...)
	return encoded
}

func encodeInteger(i int64) []byte {
	size := strconv.FormatInt(i, 10)
	encoded := append([]byte{':'}, size...)
	encoded = append(encoded, crlf...)
	return encoded
}

func encodeUnsignedInteger(i uint64) []byte {
	size := strconv.FormatUint(i, 10)

	encoded := append([]byte{':'}, size...)
	encoded = append(encoded, crlf...)
	return encoded
}

func encodeSimpleString(s string) []byte {
	encoded := append([]byte{'+'}, s...)
	encoded = append(encoded, crlf...)
	return encoded
}

func encodeBulkString(b []byte) []byte {
	size := strconv.Itoa(len(b))
	encoded := append([]byte{'$'}, size...)
	encoded = append(encoded, crlf...)
	encoded = append(encoded, b...)
	encoded = append(encoded, crlf...)
	return encoded
}

func encodeArray(a []interface{}) ([]byte, error) {
	size := strconv.Itoa(len(a))
	encoded := append([]byte{'*'}, size...)
	encoded = append(encoded, crlf...)
	for _, element := range a {
		encodedElement, err := Marshal(element)
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, encodedElement...)
	}

	return encoded, nil
}

func Marshal(v interface{}) ([]byte, error) {
	switch v.(type) {
	case nil:
		return encodeNil(), nil
	case int:
		return encodeInteger(int64(v.(int))), nil
	case int8:
		return encodeInteger(int64(v.(int8))), nil
	case int16:
		return encodeInteger(int64(v.(int16))), nil
	case int32:
		return encodeInteger(int64(v.(int32))), nil
	case int64:
		return encodeInteger(v.(int64)), nil
	case uint8:
		return encodeUnsignedInteger(uint64(v.(uint8))), nil
	case uint16:
		return encodeUnsignedInteger(uint64(v.(uint16))), nil
	case uint32:
		return encodeUnsignedInteger(uint64(v.(uint32))), nil
	case uint64:
		return encodeUnsignedInteger(v.(uint64)), nil
	case string:
		return encodeSimpleString(v.(string)), nil
	case []byte:
		return encodeBulkString(v.([]byte)), nil
	case error:
		return encodeError(v.(error)), nil
	case []interface{}:
		return encodeArray(v.([]interface{}))
	default:
		return nil, ErrUnsupportedType
	}
}
