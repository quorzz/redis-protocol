package protocol

import (
	// "bufio"
	"reflect"
	"strings"
	"testing"
)

var errorTestCase string = "-Failed..\r\n"
var statusTestCase string = "+OK\r\n"
var integerTestCase string = ":123456789\r\n"
var bulkTestCase string = "$14\r\n123456789-nice\r\n"
var mutliTestCase string = "*4\r\n$4\r\nname\r\n$5\r\nwanGe\r\n$3\r\nage\r\n$4\r\n25.0\r\n"

func TestError(t *testing.T) {
	r := parseMessage(t, errorTestCase)
	if r.Type != MessageError || r.Error.Error() != "Failed.." || !r.HasError() {
		t.Fail()
	}

	if ok, err := r.Bool(); err != nil {
		t.Error(err)
	} else if ok {
		t.Log("test error fail.")
		t.Fail()
	}
}

func TestStatus(t *testing.T) {
	r := parseMessage(t, statusTestCase)

	if r.Type != MessageStatus || string(r.Status) != "OK" {
		t.Fail()
	}

	if str, err := r.String(); err != nil {
		t.Error(err)
	} else if str != "OK" {
		t.Logf("TestStatus Fail, Expect : %s,  Got : %s", "OK", str)
	}

	if status, err := r.Bool(); err != nil {
		t.Error(err)
	} else if !status {
		t.Fail()
	}

	if i, err := r.Int64(); err != nil {
		t.Error(err)
	} else if i != 1 {
		t.Fail()
	}
}

func TestInteger(t *testing.T) {
	r := parseMessage(t, integerTestCase)

	if r.Type != MessageInt || r.Integer != 123456789 {
		t.Fail()
	}

	if status, err := r.Bool(); err != nil {
		t.Error(err)
	} else if !status {
		t.Fail()
	}

	if i, err := r.Int64(); err != nil {
		t.Error(err)
	} else if i != 123456789 {
		t.Fail()
	}

	if _, err := r.String(); err == nil {
		t.Fail()
	}
}

func TestBulk(t *testing.T) {
	r := parseMessage(t, bulkTestCase)

	if r.Type != MessageBulk {
		t.Fail()
	}

	respect := []byte("123456789-nice")
	if !reflect.DeepEqual(r.Bulk, respect) {
		t.Logf("testBulk Fail, respect %s, Got:", respect, r.Bulk)
		t.Fail()
	}

	if str, err := r.String(); err != nil {
		t.Error(err)
	} else if str != "123456789-nice" {
		t.Fail()
	}
}

func TestMutli(t *testing.T) {
	r := parseMessage(t, mutliTestCase)

	if r.Type != MessageMutli {
		t.Fail()
	}

	expect := []string{
		"name", "wanGe", "age", "25.0",
	}

	strings, err := r.Strings()
	if err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(strings, expect) {
		t.Logf("testMulti Fail : expect : %s, Got : %s", expect, strings)
	}
}
func parseMessage(t *testing.T, s string) *Message {
	// bio := bufio.NewReader(strings.NewReader(s))
	// r, e := UnpackFromReader(bio)
	r, e := NewReader(strings.NewReader(s)).ReadMessage()
	if e != nil {
		t.Error(e)
		t.Logf("Fail: %s", s)
	}
	return r
}
