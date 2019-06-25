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