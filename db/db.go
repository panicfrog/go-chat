package db

import (
	"github.com/jinzhu/gorm"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB

var diverName = "mysql"
var dataSource = "chat:3851123yw@(127.0.0.1:3306)/go_chat?charset=utf8&parseTime=True&loc=Local"

func init() {
	var err error
	DB, err = gorm.Open(diverName, dataSource)
	if err != nil {
		log.Fatal(err)
	}
	DB.AutoMigrate(&User{}, &Message{}, &Room{}, &Friend{})
	DB.LogMode(true)
}
