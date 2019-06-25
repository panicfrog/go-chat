package main

import (
	"github.com/gin-gonic/gin"
	"gochat/api"
)

func main() {
	r := gin.Default()
	api.HandleV1(r)
	r.Run(":8080")

}
