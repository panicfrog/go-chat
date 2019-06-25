package controller

import (
	"github.com/gin-gonic/gin"
	"gochat/api/service"
)

func HandleRoom(r *gin.RouterGroup) {
	var g = r.Group("/room")
	authGroup := g.Group("/auth")
	authGroup.Use(service.AuthMiddleware())
	authGroup.POST("/create", service.CreateRoom)
	authGroup.POST("/getInfo", service.GetRoomInfo)
	authGroup.POST("/addMembers", service.AddMemberToToom)
	authGroup.POST("/deleMembers", service.DeleteMemeberFromRoom)
	authGroup.POST("/empowerManagers", service.EmpowerManager)
	authGroup.POST("/callOfManagers", service.CallOffManager)
	authGroup.POST("/transferOwner", service.TransferOwner)
}