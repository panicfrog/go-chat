package chat

import (
	"github.com/gorilla/websocket"
	"sync"
)

type SocketMultipleMap struct {
	m *sync.Map
}

func NewMultipleMap() *SocketMultipleMap {
	return &SocketMultipleMap{
		m: new(sync.Map),
	}
}

func (this *SocketMultipleMap) LoadOrStore(f string, s *websocket.Conn, v *Socket) (*Socket, bool) {
	_v, loaded := this.m.LoadOrStore(f, v)
	if loaded {
		return _v.(*Socket), loaded
	} else {
		this.m.Store(s, v)
		this.m.Store(f, v)
		return v, loaded
	}
}

func (this *SocketMultipleMap) Store(f string, s *websocket.Conn, v *Socket) {
	this.m.Store(s, v)
	this.m.Store(f, v)
}

func (this *SocketMultipleMap) LoadWithId(f string) (s *Socket, exited bool)  {
	v, exited := this.m.Load(f)
	if exited {
		s = v.(*Socket)
	}
	return
}

func (this *SocketMultipleMap) LoadWithConn(s *websocket.Conn) (socket *Socket, exited bool) {
	v, exited := this.m.Load(s)
	if exited {
		socket = v.(*Socket)
	}
	return
}

func (this *SocketMultipleMap) DeleteWithId(f string)  {
	if v, exited := this.LoadWithId(f); exited {
		this.m.Delete(v.Conn)
		this.m.Delete(v.Id)
	}
}

func (this *SocketMultipleMap) DeleteWithConn(s *websocket.Conn) {
	if v, exited := this.LoadWithConn(s); exited {
		this.m.Delete(v.Conn)
		this.m.Delete(v.Id)
	}
}