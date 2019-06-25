package controller

import (
	"github.com/gin-gonic/gin"
	"gochat/api/service"
)

func HandleUser(r *gin.RouterGroup) {
	var g = r.Group("/user")
	g.POST("/register", service.Register)
	g.POST("login", service.Login)
	authGroup := g.Group("/auth")
	authGroup.Use(service.AuthMiddleware())

	authGroup.POST("/friends", service.Friends)
	authGroup.POST("/addFriends", service.AddFriends)
	authGroup.POST("/removeFriends", service.RemoveFriends)
	authGroup.POST("/rooms", service.Rooms)
}