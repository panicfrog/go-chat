package chat

import (
	"encoding/json"
	"gochat/db"
	"log"
	"sort"
)


func encodeMessageId(first string, second string) (string) {
	ss := []string{first, second}
	sort.Strings(ss)
	v, err := json.Marshal(ss)
	if err != nil {
		log.Println(err)
	}
	return string(v)
}

func decodeMessageId(v string) (first string, second string) {
	ss := []string{}
	if err := json.Unmarshal([]byte(v), &ss); err != nil {
		log.Println(err)
	} else if len(ss) == 2 {
		first = ss[0]
		second = ss[1]
	}
	return
}

func SendMessage(message db.Message) {

	msg, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return
	}
	// 单聊
	if message.MessageType == db.MessageTypeSigleChat {
		SendToUser(message.From, message.To, string(msg))

	}

	// 群聊
	 if message.MessageType == db.MessageTypeGroupChat {
		SendToRoom(message.From, message.To, string(msg))
	 }
}


func SendToRoom(from string, roomId string, content string) {
	socket, exited := SocketBucket.LoadWithId(from)
	if !exited {
		log.Println("用户'", from, "'不存在")
		return
	}
	room := db.Room{}
	err := db.DB.Where(&db.Room{RoomId:roomId}).First(&room).Error
	if err != nil {
		log.Println(err)
		return
	}

	users := []db.User{}
	err = db.DB.Model(&room).Association("Users").Find(&users).Error

	var isInRoom = false
	for _, u := range users {
		if u.UserName == from {
			isInRoom = true
			break
		}
	}
	if !isInRoom {
		err := socket.SendMessage([]byte("只有房间的成员才能向房间发送消息"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	if err != nil {
		log.Println(err)
		err := socket.SendMessage([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
		return
	}

	for _, u := range users {
		s, exited := SocketBucket.LoadWithId(u.UserName)
		if !exited {
			continue
		}
		err = s.SendMessage([]byte(content))
		if err != nil {
			log.Println(err)
		}
	}

	// TODO: 将消息存到数据库中
}

func SendToUser(from string, to string, message string,) {

	s, exited := SocketBucket.LoadWithId(from)
	if !exited {
		println("'", from, "'不在线")
		return
	}

	user := db.User{}
	dbErr := db.DB.Where(&db.User{UserName:from}).First(&user).Error
	if dbErr != nil {
		log.Println(dbErr)
		return
	}

	friends := []db.User{}
	dbErr = db.DB.Model(&user).Association("Friends").Find(&friends).Error
	if dbErr != nil {
		log.Println(dbErr)
		return
	}

	isFriend := false
	for _,u := range friends {
		if u.UserName == to {
			isFriend = true
			break
		}
	}
	if !isFriend {
		if err := s.SendMessage([]byte("只能给好友发送消息")); err != nil {
			println(err)
		}
		return
	}
	toS, exited := SocketBucket.LoadWithId(to)
	if !exited {
		return
	}
	err := toS.SendMessage([]byte(message))
	if err != nil {
		log.Println(err)
	}
	// TODO: 将消息存到数据库中 发送之后将to转成数据库保存的方式(使用「encodeMessageId」和「decodeMessageId」)
}