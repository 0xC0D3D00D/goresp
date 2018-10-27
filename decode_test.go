package goresp

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

func BenchmarkReadSimpleString(b *testing.B) {
	reader := bytes.NewReader([]byte{'a', 'b', 'c', '\r', '\n'})
	for n := 0; n < b.N; n++ {
		reader.Seek(0, io.SeekStart)
		readSimpleString(reader)
	}
}

func BenchmarkReadInteger(b *testing.B) {
	reader := bytes.NewReader([]byte{'1', '2', '3', '4', '\r', '\n'})
	for n := 0; n < b.N; n++ {
		reader.Seek(0, io.SeekStart)
		readInteger(reader)
	}
}

func BenchmarkReadBulkString(b *testing.B) {
	reader := bytes.NewReader([]byte{'3', '\r', '\n', 'a', 'b', 'c', '\r', '\n'})
	for n := 0; n < b.N; n++ {
		reader.Seek(0, io.SeekStart)
		readBulkString(reader)
	}
}

func BenchmarkReadArray(b *testing.B) {
	reader := bytes.NewReader([]byte{'2', '\r', '\n', '$', '3', '\r', '\n', 'f', 'o', 'o', '\r', '\n', '$', '3', '\r', '\n', 'b', 'a', 'r', '\r', '\n'})
	for n := 0; n < b.N; n++ {
		reader.Seek(0, io.SeekStart)
		readArray(reader)
	}
}

func TestRead(t *testing.T) {
	testCases := []struct {
		msg  []byte
		resp interface{}
		err  error
	}{
		{
			[]byte{'\r', '\n'},
			nil,
			nil,
		}, // empty
		{
			[]byte{'*', '2', '\r', '\n',
				'$', '3', '\r', '\n',
				'f', 'o', 'o', '\r', '\n',
				'$', '3', '\r', '\n',
				'b', 'a', 'r', '\r', '\n'},
			[]interface{}{
				[]byte{'f', 'o', 'o'},
				[]byte{'b', 'a', 'r'},
			},
			nil,
		}, // array of two bulk strings
		{
			[]byte{':', '1', '\r', '\n'},
			int64(1),
			nil,
		}, // an integer
		{
			[]byte{'*', '5', '\r', '\n',
				':', '1', '\r', '\n',
				':', '2', '\r', '\n',
				':', '3', '\r', '\n',
				':', '4', '\r', '\n',
				'$', '6', '\r', '\n',
				'f', 'o', 'o', 'b', 'a', 'r', '\r', '\n',
			},
			[]interface{}{
				int64(1),
				int64(2),
				int64(3),
				int64(4),
				[]byte{'f', 'o', 'o', 'b', 'a', 'r'},
			},
			nil,
		}, // mixed array
		{
			[]byte{'+', 'O', 'K', '\r', '\n'},
			[]byte{'O', 'K'},
			nil,
		}, // simple string
		{
			[]byte{'-', 'E', 'R', 'R', '\r', '\n'},
			errors.New("ERR"),
			nil,
		}, // error
		{
			[]byte{'\r'},
			nil,
			ErrUnexpectedEOF,
		}, // empty (bad termination)
		{
			[]byte{'!', 'I', 'N', 'V', 'A', 'L', 'I', 'D', '\r', '\n'},
			nil,
			ErrInvalidMessage,
		}, // invalid type
	}

	for _, testCase := range testCases {
		reader := bytes.NewReader(testCase.msg)
		resp, err := read(bufio.NewReader(reader))
		if !reflect.DeepEqual(resp, testCase.resp) || err != testCase.err {
			t.Fatalf("Case %v:\nExpected resp=%v and err=%v, Actual resp=%v, err=%v", testCase.msg, testCase.resp, testCase.err, resp, err)
		}
	}
}

func TestReadArray(t *testing.T) {
	testCases := []struct {
		msg  []byte
		resp interface{}
		err  error
	}{
		{
			[]byte{},
			nil,
			ErrUnexpectedEOF,
		}, // no input
		{
			[]byte{'0', '\r', '\n'},
			[]interface{}{},
			nil,
		}, // empty
		{
			[]byte{'2', '\r', '\n',
				'$', '3', '\r', '\n',
				'f', 'o', 'o', '\r', '\n',
				'$', '3', '\r', '\n',
				'b', 'a', 'r', '\r', '\n'},
			[]interface{}{
				[]byte{'f', 'o', 'o'},
				[]byte{'b', 'a', 'r'},
			},
			nil,
		}, // two bulk strings
		{
			[]byte{'3', '\r', '\n', ':', '1', '\r', '\n', ':', '2', '\r', '\n', ':', '3', '\r', '\n'},
			[]interface{}{int64(1), int64(2), int64(3)},
			nil,
		}, // three integers
		{
			[]byte{'5', '\r', '\n',
				':', '1', '\r', '\n',
				':', '2', '\r', '\n',
				':', '3', '\r', '\n',
				':', '4', '\r', '\n',
				'$', '6', '\r', '\n',
				'f', 'o', 'o', 'b', 'a', 'r', '\r', '\n',
			},
			[]interface{}{
				int64(1),
				int64(2),
				int64(3),
				int64(4),
				[]byte{'f', 'o', 'o', 'b', 'a', 'r'},
			},
			nil,
		}, // mixed array
		{
			[]byte{'3', '\r', '\n', ':', '1', '\r', '\n', ':', '2', '\r', '\n'},
			nil,
			ErrUnexpectedEOF,
		}, // array of integers (bad size)
	}

	for _, testCase := range testCases {
		reader := bytes.NewReader(testCase.msg)
		resp, err := readArray(bufio.NewReader(reader))
		if !reflect.DeepEqual(resp, testCase.resp) || err != testCase.err {
			t.Fatalf("Case %v:\nExpected resp=%v and err=%v, Actual resp=%v, err=%v", testCase.msg, testCase.resp, testCase.err, resp, err)
		}
	}
}

func TestReadBulkString(t *testing.T) {
	testCases := []struct {
		msg  []byte
		resp []byte
		err  error
	}{
		{
			[]byte{},
			nil,
			ErrUnexpectedEOF,
		}, // no input
		{
			[]byte{'0', '\r', '\n', '\r', '\n'},
			[]byte{},
			nil,
		}, // empty string
		{
			[]byte{'3', '\r', '\n', 'a', 'b', 'c', '\r', '\n'},
			[]byte{'a', 'b', 'c'},
			nil,
		}, // valid
		{
			[]byte{'3', '\r', '\n', 'a', 'b', 'c'},
			nil,
			ErrUnexpectedEOF,
		}, // valid (no termination)
		{
			[]byte{'3', '\r', '\n', 'a', 'b', 'c', '\r', '\r'},
			nil,
			ErrInvalidMessage,
		}, // valid (bad termination)
		{
			[]byte{'3', '\r', '\n'},
			nil,
			ErrUnexpectedEOF,
		}, // Unexpected termination
		{
			[]byte{'2', '\r', '\n', '\r', '\n', '\r', '\n'},
			[]byte{'\r', '\n'},
			nil,
		}, // binary safe check
		{
			[]byte{'0', '\r', '\n'},
			nil,
			ErrUnexpectedEOF,
		}, // invalid empty string
		{
			[]byte{'B', 'A', 'D', '\r', '\n'},
			nil,
			ErrInvalidMessage,
		}, // bad size
	}

	for _, testCase := range testCases {
		reader := bytes.NewReader(testCase.msg)
		resp, err := readBulkString(bufio.NewReader(reader))
		if !reflect.DeepEqual(resp, testCase.resp) || err != testCase.err {
			t.Fatalf("Case %v:\nExpected resp=%#v and err=%v, Actual resp=%#v, err=%v", testCase.msg, testCase.resp, testCase.err, resp, err)
		}
	}
}

func TestReadInteger(t *testing.T) {
	testCases := []struct {
		msg  []byte
		resp int64
		err  error
	}{
		{
			[]byte{},
			0,
			ErrUnexpectedEOF,
		}, // no input
		{
			[]byte{'\r', '\n'},
			0,
			ErrInvalidMessage,
		}, // empty
		{
			[]byte{'1', '2', '3', '4', '\r', '\n'},
			1234,
			nil,
		}, // 1234
		{
			[]byte{'1', '2', '3', '4', '\r'},
			0,
			ErrUnexpectedEOF,
		}, // 1234 (bad termination)
		{
			[]byte{'1', 'E', 'R', 'R', '\r', '\n'},
			0,
			ErrInvalidMessage,
		}, // invalid integer
	}

	for _, testCase := range testCases {
		reader := bytes.NewReader(testCase.msg)
		resp, err := readInteger(bufio.NewReader(reader))
		if resp != testCase.resp || err != testCase.err {
			t.Fatalf("Case %v:\nExpected resp=%#v and err=%v, Actual resp=%#v, err=%v", testCase.msg, testCase.resp, testCase.err, resp, err)
		}
	}
}

func TestReadSimpleString(t *testing.T) {
	testCases := []struct {
		msg  []byte
		resp []byte
		err  error
	}{
		{
			[]byte{},
			nil,
			ErrUnexpectedEOF,
		}, // no input
		{
			[]byte{'\r', '\n'},
			[]byte{},
			nil,
		}, // empty
		{
			[]byte{'a', 'b', 'c', '\r', '\n'},
			[]byte{'a', 'b', 'c'},
			nil,
		}, // abc
		{
			[]byte{'a', 'b', 'c', '\r'},
			nil,
			ErrUnexpectedEOF,
		}, // abc (bad termination)
	}

	for _, testCase := range testCases {
		reader := bytes.NewReader(testCase.msg)
		resp, err := readSimpleString(bufio.NewReader(reader))
		if !reflect.DeepEqual(resp, testCase.resp) || err != testCase.err {
			t.Fatalf("Case %v:\nExpected resp=%#v and err=%v, Actual resp=%#v, err=%v", testCase.msg, testCase.resp, testCase.err, resp, err)
		}
	}
}
