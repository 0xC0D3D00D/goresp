package goresp

import (
	"errors"
	"io"
	"strconv"
)

var (
	// ErrUnexpectedEOF thrown when unmarshaling an incomplete message
	ErrUnexpectedEOF = errors.New("Unexpected EOF")
	// ErrInvalidMessage thrown when the message is not in any form of standard RESP format
	ErrInvalidMessage = errors.New("Invalid message")

	typeError        = byte('-')
	typeSimpleString = byte('+')
	typeInteger      = byte(':')
	typeBulkString   = byte('$')
	typeArray        = byte('*')
)

func readInteger(reader io.Reader) (int64, error) {
	bytes, err := readSimpleString(reader)
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(string(bytes), 10, 64)
	if err != nil {
		return 0, ErrInvalidMessage
	}
	return i, nil
}

func readSimpleString(reader io.Reader) ([]byte, error) {
	oneByte := make([]byte, 1)
	readBytes := 0
	var err error

	str := []byte{}
	for {
		_, err = reader.Read(oneByte)
		if err == io.EOF {
			return nil, ErrUnexpectedEOF
		}
		if err != nil {
			return nil, err
		}

		if oneByte[0] == '\r' {
			_, err = reader.Read(oneByte)
			if err == io.EOF {
				return nil, ErrUnexpectedEOF
			}
			if err != nil {
				return nil, err
			}
			if oneByte[0] == '\n' {
				break
			}
		}

		readBytes++
		str = append(str, oneByte[0])
	}

	return str[:readBytes], nil
}

func readBulkString(reader io.Reader) ([]byte, error) {
	size, err := readInteger(reader)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		dummyBytes := make([]byte, 2)
		_, err = reader.Read(dummyBytes)
		if err != nil {
			if err == io.EOF {
				return nil, ErrUnexpectedEOF
			}
			return nil, err
		}
		return []byte{}, nil
	}

	str := make([]byte, size+2)
	readBytes, err := reader.Read(str)
	if err != nil {
		if err == io.EOF {
			return nil, ErrUnexpectedEOF
		}
		return nil, err
	}
	if readBytes != len(str) {
		return nil, ErrUnexpectedEOF
	}
	if str[size] != '\r' || str[size+1] != '\n' {
		return nil, ErrInvalidMessage
	}

	return str[:size], nil
}

func readArray(reader io.Reader) (interface{}, error) {
	size, err := readInteger(reader)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return []interface{}{}, nil
	}

	array := make([]interface{}, size)
	for i := int64(0); i < size; i++ {
		element, err := read(reader)
		if err != nil {
			if err == io.EOF {
				return nil, ErrUnexpectedEOF
			}
			return nil, err
		}
		array[i] = element
	}

	return array, nil
}

func read(reader io.Reader) (interface{}, error) {
	typeIndicator := []byte{0}
	_, err := reader.Read(typeIndicator)
	if err != nil {
		return nil, err
	}
	switch typeIndicator[0] {
	case '\r': //empty
		reader.Read(typeIndicator) // CHECK
		if typeIndicator[0] == '\n' {
			return nil, nil
		}
		return nil, ErrUnexpectedEOF
	case typeSimpleString:
		return readSimpleString(reader)
	case typeError:
		errMsg, err := readSimpleString(reader)
		if errMsg == nil {
			return nil, err
		}
		return errors.New(string(errMsg)), err
	case typeInteger:
		return readInteger(reader)
	case typeBulkString:
		return readBulkString(reader)
	case typeArray:
		return readArray(reader)
	default:
		return nil, ErrInvalidMessage
	}
}

// Unmarshal decodes RESP messages to its corresponding golang type
func Unmarshal(reader io.Reader) (interface{}, error) {
	return read(reader)
}
