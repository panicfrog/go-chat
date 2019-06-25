package api

import (
	"github.com/gin-gonic/gin"
	"gochat/api/controller"
	"gochat/chat"
)

func HandleV1(r *gin.Engine)  {
	v1Group := r.Group("v1")
	controller.CommonHandle(v1Group)
	controller.HandleUser(v1Group)
	controller.HandleRoom(v1Group)
	chat.HanderWebsocket(v1Group)
}