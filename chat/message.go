package chat

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
	return base64.StdEncoding.EncodeToString(v)
}

func decodeMessageId(v string) (first string, second string) {
	b, e := base64.StdEncoding.DecodeString(v)
	if e != nil {
		log.Println(e)
		return
	}
	ss := []string{}
	if err := json.Unmarshal(b, &ss); err != nil {
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

	var e error
	// 单聊
	if message.MessageType == db.MessageTypeSigleChat {
		e = sendToUser(message.From, message.To, string(msg))
	}

	// 群聊
	 if message.MessageType == db.MessageTypeGroupChat {
		e =	sendToRoom(message.From, message.To, string(msg))
	 }

	if e != nil {
		return
	}

	// 单聊时，将to字段统一处理
	if message.MessageType == db.MessageTypeSigleChat {
		message.To = encodeMessageId(message.From, message.To)
	}
 	if	dbErr := db.DB.Create(&message).Error; dbErr != nil {
 		log.Println(dbErr)
	}
}


func sendToRoom(from string, roomId string, content string) (error) {
	socket, exited := SocketBucket.LoadWithId(from)
	if !exited {
		log.Println("用户'", from, "'不存在")
		return errors.New("用户不存在")
	}
	room := db.Room{}
	err := db.DB.Where(&db.Room{RoomId:roomId}).First(&room).Error
	if err != nil {
		log.Println(err)
		return err
	}

	users := []db.User{}
	err = db.DB.Model(&room).Association("Users").Find(&users).Error

	if err != nil {
		log.Println(err)
		_err := socket.SendMessage([]byte(err.Error()))
		if _err != nil {
			log.Println(_err)
		}
		return err
	}

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
		return err
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
	var e error
	return  e
}

func sendToUser(from string, to string, message string,) (error) {

	s, exited := SocketBucket.LoadWithId(from)
	if !exited {
		println("'", from, "'不在线")
		return errors.New("不在线")
	}

	user := db.User{}
	dbErr := db.DB.Where(&db.User{UserName:from}).First(&user).Error
	if dbErr != nil {
		log.Println(dbErr)
		return dbErr
	}

	friends := []db.User{}
	dbErr = db.DB.Model(&user).Association("Friends").Find(&friends).Error
	if dbErr != nil {
		log.Println(dbErr)
		return dbErr
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
		return errors.New("只能给好友发送消息")
	}
	toS, exited := SocketBucket.LoadWithId(to)
	if !exited {
		return errors.New("用户不存在")
	}
	err := toS.SendMessage([]byte(message))
	if err != nil {
		log.Println(err)
	}
	var e error
	return  e
}