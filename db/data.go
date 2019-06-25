package db

import (
	"github.com/jinzhu/gorm"
)

/*
	message : message_id 消息id, form: 发消息的人, 消息类型: 单聊，群聊, to: 接受方（房间号，或者用户ID拼接的）
	拼接ID的方式 username->升序排序之后想加->base64
*/


type MessageType int
const (
	MessageTypeSigleChat = iota
	MessageTypeGroupChat
)

type User struct {
	gorm.Model				`json:"-"`
	UserName string 		`gorm:"type:varchar(50);not null;unique" json:"user_name"`
	Passwd string 			`gorm:"type:varchar(50);not null" json:"-"`
	Friends []User			`gorm:"many2many:friendship;association_jointable_foreignkey:friend_id" json:"-"`
	Rooms []Room			`gorm:"many2many:room_users" json:"-"`
}

type Message struct {
	gorm.Model				`json:"-"`
	MessageId string  		`gorm:"type:varchar(100);not null" json:"message_id"`
	MessageType MessageType `json:"message_type" json:"message_type"`
	From string 			`gorm:"type:varchar(50);not null" json:"from"`
	To string				`gorm:"type:varchar(100);not null" json:"to"`
	Content string 			`gorm:"type:varchar(500);not null" json:"content"`
}

type Room struct {
	gorm.Model				`json:"-"`
	Owned string 			`gorm:"type:varchar(50);not null" json:"owned"`
	RoomId string 			`gorm:"type:varchar(100);not null;unique" json:"room_id"`
	Managers []User			`gorm:"many2many:room_managers" json:"mangaers"`
	Users []User			`gorm:"many2many:room_users" json:"users"`
}

type Friend struct {
	gorm.Model				`json:"-"`
	User string				`gorm:"type:varchar(50);not null" json:"user"`
	Friend string			`gorm:"type:varchar(50);not null" json:"friend"`
}