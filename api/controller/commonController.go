package controller

import (
	"github.com/gin-gonic/gin"
	"gochat/api/service"
)

func CommonHandle(r *gin.RouterGroup) {
	g := r.Group("common")
	g.GET("/hello", service.CommSuccess)
	g.GET("/fail", service.CommFail)
	g.GET("/testDB", service.CommDB)

	authGroup := g.Group("/auth")
	authGroup.Use(service.AuthMiddleware())

	authGroup.GET("/user", service.CommAuth)
}

