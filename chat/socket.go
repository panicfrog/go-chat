package chat

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"gochat/db"
	"gochat/variable"
	"log"
	"sync"
	"time"
)

type MessageSender interface {
	SendMessage(data []byte) ( err error)
	LoadWithId(f string) (s *Socket, exited bool)
}

type Socket struct {
	sender 	chan []byte
	reciver chan []byte
	closer 	chan byte
	mutex sync.Mutex
	isClose bool
	Id string
	Platform variable.Platform
	Conn *websocket.Conn
}

func NewSocket(id string, platform variable.Platform, conn *websocket.Conn) *Socket{
	return &Socket{
		sender:	 make(chan []byte, 1000),
		reciver: make(chan []byte, 1000),
		closer:	 make(chan byte, 1),
		mutex: sync.Mutex{},
		isClose: false,
		Id: id,
		Platform: platform,
		Conn: conn,
	}
}

func (this *Socket) Close() {
	this.mutex.Lock()
	if !this.isClose {
		err := this.Conn.Close()
		log.Println(err)
		close(this.closer)
		this.isClose = true
	}
	this.mutex.Unlock()
}

// 运行
func (this *Socket) RunLoop() {
	go this.sendRoop()
	go this.reciveRoop()
	go this.DealMessage()
	this.onClose()
}

// 处理业务
func (this *Socket) DealMessage() {
	var data []byte
	var err error
	for {
		if data, err = this.ReciveMessage(); err != nil {
			goto ERR
		}

		msg := db.Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Println(err)
		}
		if msg.From != this.Id {
			_ = this.SendMessage([]byte("请确定是你发出的消息"))
		} else {
			SendMessage(msg)
		}
	}
	ERR:
		this.Close()
}

// 接受消息
func (this *Socket) ReciveMessage() (data []byte, err error) {
	select {
		case data = <- this.reciver:
		case <- this.closer:
				err = errors.New("连接关闭")
	}
	return
}

// 发送消息
func (this *Socket) SendMessage(data []byte) ( err error) {
	select {
		case this.sender <- data:
		case <- this.closer:
				err = errors.New("连接关闭")
	}
	return
}

// 接受loop
func (this *Socket) reciveRoop() {
	var	data []byte
	var err error
	for {
		if _, data, err = this.Conn.ReadMessage(); err != nil {
			goto ERR
		}
		select {
			case this.reciver <- data:
			case <- this.closer:
				goto ERR
		}
	}
	ERR:
		this.Close()
}

// 发送loop
func (this *Socket) sendRoop() {
	var data []byte
	var err error
	for {
		select {
			case data = <- this.sender:
			case <- this.closer:
				goto ERR
		}
		if err = this.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println(err)
			goto ERR
		}
	}
	ERR:
		this.Close()
}


func (this *Socket) onClose()  {
	this.Conn.SetCloseHandler(func (code int, text string) error  {
		message := websocket.FormatCloseMessage(code, "")
		_ = this.Conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))
		log.Println("关闭连接")
		this.Close()
		return nil
	})
}