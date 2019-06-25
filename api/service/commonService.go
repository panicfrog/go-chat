package service

import (
	"github.com/gin-gonic/gin"
	"gochat/db"
	"gochat/variable"
	"log"
)

func CommSuccess(c *gin.Context) {
	sendSuccess("请求成功", "hello world", c)
}

func CommFail(c *gin.Context) {
	sendFail("请求失败", c)
}

func CommDB(c *gin.Context) {
	users := []db.User{}
	db.DB.Find(&users)
	//_, err := db.DBengine.Query("select * from user")
	//if err != nil {
	//	log.Println(err)
	//}
}

func CommAuth(c *gin.Context) {
	userName, exit := c.Get(variable.AuthKey)
	if !exit {
		log.Println("没有获取到用户信息")
		sendServerInternelError("没有获取到用户信息", c)
		return
	}
	sendSuccess("获取到用户数据", userName, c)
}