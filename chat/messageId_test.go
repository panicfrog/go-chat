package chat

import (
	"fmt"
	"testing"
)

func TestMessageId(t *testing.T) {
	first, second := "xiaoming", "xiaohuang"
	v := encodeMessageId(first, second)
	_first, _second := decodeMessageId(v)
	if first == _first {
		if second != _second {
			e := fmt.Sprintf("first: %s, second: %s, _first: %s, _second: %s", first, second, _first, _second)
			t.Error(e)
		}
	} else if first == _second {
		if second != _first {
			e := fmt.Sprintf("first: %s, second: %s, _first: %s, _second: %s", first, second, _first, _second)
			t.Error(e)
		}
	} else {
		e := fmt.Sprintf("first: %s, second: %s, _first: %s, _second: %s", first, second, _first, _second)
		t.Error(e)
	}
}

func TestMessageIDDecode(t *testing.T) {
	f, s := decodeMessageId("WyJ4aWFvbGlhbmciLCJ4aWFvbml1Il0=")
	if f == "xiaoniu" {
		if s != "xiaoliang" {
			e := fmt.Sprintf("f: %s, s: %s", f, s)
			t.Error(e)
		}
	} else if f == "xiaoliang" {
		if s != "xiaoniu" {
			e := fmt.Sprintf("f: %s, s: %s", f, s)
			t.Error(e)
		}
	} else {
		e := fmt.Sprintf("f: %s, s: %s", f, s)
		t.Error(e)
	}
}