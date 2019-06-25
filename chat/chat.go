package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"gochat/variable"
	"net/http"
	"sync"
)

var SocketBucket *SocketMultipleMap
var UnAuthSocketBucket *sync.Map

func init() {
	SocketBucket = NewMultipleMap()
	UnAuthSocketBucket = new(sync.Map)
}

func HanderWebsocket(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	authGroup.Use(ChatAuthMiddleware())
	authGroup.GET("/ws", connetWebsocet)
}

func connetWebsocet(c *gin.Context){

	_id, exited := c.Get(variable.AuthKey)
	id := ""
	if !exited {
		log.Println("未授权的websocket连接")
	} else {
		id = _id.(string)
	}

	platform, exited := c.Get(variable.PlatformKey)
	p := variable.PlatformUnknow
	if exited {
		p = platform.(variable.Platform)
	}

	verify := false
	if p == variable.PlatformUnknow && id != "" {
		verify = true
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:1024,
		WriteBufferSize:1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: verify,
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	socket := NewSocket(id, p, conn)
	socket.RunLoop()
	// 从map中删除对应的连接
	go func() {
		for {
			select {
			case <- socket.closer:
				goto CLOSE
			default:
			}
		}
	CLOSE:
		SocketBucket.DeleteWithConn(socket.Conn)
		UnAuthSocketBucket.Delete(socket.Conn)
		log.Println("已经删除")
	}()

	if verify {
		_, loaded := SocketBucket.LoadOrStore(id, conn, socket)
		if loaded {
			log.Println("该id已经绑定对应的socket")
			socket.Close()
			return
		}
	} else {
		_, loaded := UnAuthSocketBucket.LoadOrStore(conn, socket)
		if loaded {
			log.Println("该socket已在管理中")
			return
		}
	}


}