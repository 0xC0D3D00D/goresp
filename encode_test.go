package goresp

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func BenchmarkEncodeInteger(b *testing.B) {
	for n := 0; n < b.N; n++ {
		encodeInteger(1234)
	}
}

func BenchmarkEncodeSimpleString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		encodeSimpleString("OK")
	}
}

func makeEncodeBulkStringBench(b *testing.B, size int) func(b *testing.B) {
	bulk := make([]byte, size)
	n, err := rand.Read(bulk)
	if err != nil || n != size {
		b.Fatal(err)
	}
	return func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			encodeBulkString(bulk)
		}
	}
}

func BenchmarkEncodeBulkString(b *testing.B) {
	for i := 1; i <= (1 << 20); i = i * 32 {
		b.Run(fmt.Sprintf("%dBytes", i), makeEncodeBulkStringBench(b, i))
	}
}

func TestEncodeNil(t *testing.T) {
	expected := []byte{'$', '-', '1', '\r', '\n'}
	actual := encodeNil()

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestEncodeUnsignedInteger(t *testing.T) {
	expected := []byte{':', '4', '2', '9', '4', '9', '6', '7', '2', '9', '7', '\r', '\n'}
	actual := encodeUnsignedInteger(4294967297) // 2^32+1

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestEncodeInteger(t *testing.T) {
	expected := []byte{':', '-', '8', '\r', '\n'}
	actual := encodeInteger(-8)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestEncodeError(t *testing.T) {
	expected := []byte{'-', 'E', 'R', 'R', '\r', '\n'}
	actual := encodeError(errors.New("ERR"))

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestEncodeSimpleString(t *testing.T) {
	expected := []byte{'+', 'O', 'K', '\r', '\n'}
	actual := encodeSimpleString("OK")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestEncodeBulkString(t *testing.T) {
	expected := []byte{'$', '2', '\r', '\n', 'O', 'K', '\r', '\n'}
	actual := encodeBulkString([]byte{'O', 'K'})

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected %v got %v", expected, actual)
	}
}

func TestEncodeArray(t *testing.T) {
	testCases := []struct {
		msg  []byte
		resp []interface{}
		err  error
	}{
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
			[]byte{'*', '5', '\r', '\n',
				':', '1', '\r', '\n',
				':', '2', '\r', '\n',
				':', '3', '\r', '\n',
				':', '4', '\r', '\n',
				'$', '6', '\r', '\n',
				'f', 'o', 'o', 'b', 'a', 'r', '\r', '\n',
			},
			[]interface{}{
				int(1),
				int8(2),
				int32(3),
				int64(4),
				[]byte{'f', 'o', 'o', 'b', 'a', 'r'},
			},
			nil,
		}, // mixed array
	}

	for _, testCase := range testCases {
		msg, err := encodeArray(testCase.resp)
		if !reflect.DeepEqual(msg, testCase.msg) || err != testCase.err {
			t.Fatalf("Case %v:\nExpected resp=%v and err=%v, Actual resp=%v, err=%v", testCase.msg, testCase.resp, testCase.err, msg, err)
		}
	}
}
