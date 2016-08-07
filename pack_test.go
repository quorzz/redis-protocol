package protocol

import (
	"reflect"
	"testing"
)

type packTestCase struct {
	value     interface{}
	respect   []byte
	errorCase bool
}

type normalizeArgsTestCase struct {
	value       interface{}
	respect     []interface{}
	compareFunc func([]interface{}, []interface{}) bool
	errorCase   bool
}

var packTestCases = []packTestCase{
	{
		"fooBar123: -/_", []byte("*1\r\n$14\r\nfooBar123: -/_\r\n"), false,
	},
	{
		000, []byte("*1\r\n$1\r\n0\r\n"), false,
	},
	{
		3.141592653, []byte("*1\r\n$11\r\n3.141592653\r\n"), false,
	},
	{
		0.628, []byte("*1\r\n$5\r\n0.628\r\n"), false,
	},
	{
		false, []byte("*1\r\n$1\r\n0\r\n"), false,
	},
	{
		true, []byte("*1\r\n$1\r\n1\r\n"), false,
	},
	{
		nil, []byte("*1\r\n$0\r\n\r\n"), false,
	},
	{
		[]byte{'H', 'i', '!'}, []byte("*1\r\n$3\r\nHi!\r\n"), false,
	},

	// error case
	{
		[]string{"nice", "to", "meet", "you"}, nil, true,
	},
	{
		map[string]int{"n": 9, "o": 7}, nil, true,
	},
}

var normalizeArgsTestCases = []normalizeArgsTestCase{
	{
		1234, []interface{}{1234}, sliceCompare, false,
	},
	{
		[]string{"3.14", "nice"}, []interface{}{"3.14", "nice"}, sliceCompare, false,
	},
	{
		map[int]string{90: "1911", 12: "ce"}, []interface{}{90, "1911", 12, "ce"}, mapCompare, false,
	},
	{
		map[interface{}]string{1: "2", "3": "4"}, []interface{}{"3", "4", 1, "2"}, mapCompare, false,
	},
	{
		map[interface{}]interface{}{"abc": map[int]string{123: "ppp"}},
		[]interface{}{
			"abc", map[int]string{123: "ppp"},
		},
		mapCompare,
		false,
	},

	// error case
	{
		[]int{1234, 8901}, []interface{}{"1234", 8901}, sliceCompare, true,
	},
	{
		map[string]string{"a": "b", "x": "y"},
		[]interface{}{"a", "b", "y", "x"},
		mapCompare,
		true,
	},
}

func TestPack(t *testing.T) {

	for k, testCase := range packTestCases {

		if packBytes, e := PackCommand(testCase.value); e != nil && !testCase.errorCase {
			t.Error(e)
		} else if !reflect.DeepEqual(packBytes, testCase.respect) {

			t.Logf("the %dth pack testCase Failed\r\n respect : %s \r\n Got : %s\r\n", k, string(testCase.respect), string(packBytes))
			t.Fail()
		}
	}
}

func TestNormalizeArgs(t *testing.T) {

	for k, testCase := range normalizeArgsTestCases {

		if normalizeArgs := NormalizeArgs(testCase.value); !testCase.compareFunc(testCase.respect, normalizeArgs) && !testCase.errorCase {

			t.Logf("the %dth pack normalizeArgs Failed\r\n respect : %v \r\n Got : %v\r\n", k, testCase.respect, normalizeArgs)
			t.Fail()
		}
	}
}

func mapCompare(a []interface{}, b []interface{}) bool {
	if len(a) != len(b) || len(a)&1 == 1 {
		return false
	}
	for k, v := range a {
		for kk, vv := range b {
			if k&1 == 0 && kk&1 == 0 && reflect.DeepEqual(v, vv) && reflect.DeepEqual(a[k+1], b[kk+1]) {
				a = delSlice(a, k)
				a = delSlice(a, k)
				b = delSlice(b, kk)
				b = delSlice(b, kk)
				return mapCompare(a, b)
			}
		}
	}
	return len(a) == 0 && len(a) == len(b)
}

func sliceCompare(a []interface{}, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		for kk, vv := range b {
			if reflect.DeepEqual(v, vv) {
				a = delSlice(a, k)
				b = delSlice(b, kk)
				return sliceCompare(a, b)
			}
		}
	}
	return len(a) == 0 && len(b) == 0
}

func delSlice(v []interface{}, k int) []interface{} {
	if k+1 == len(v) {
		return v[:k]
	} else {
		s := make([]interface{}, len(v))
		copy(s, v)
		return append(s[:k], s[k+1:]...)
	}
}
