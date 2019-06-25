package chat

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gochat/variable"
	"log"
	"testing"
)

func TestNewMultipleMap(t *testing.T) {
	m := NewMultipleMap()
	id := "abcdefghight"
	conn := &websocket.Conn{}
	socket := NewSocket(id , variable.PlatformIOS, conn)
	m.Store(id, conn, socket)
	s,e := m.LoadWithId(id)
	st := fmt.Sprintf("exited: %t get: %p, source: %p, get id: %s, source id: %s", e, s.Conn, socket.Conn, s.Id, socket.Id)
	if !e || s.Id != socket.Id {
		t.Error(errors.New(st))
	} else {
		log.Println(st)
	}
}