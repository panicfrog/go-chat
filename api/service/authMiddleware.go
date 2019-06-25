package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gochat/toolUtils"
	"gochat/variable"
	"log"
	"net/http"
)


// TODO: 还有平台属性没有处理

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			log.Println("无token")
			sendHTTPError(c, http.StatusUnauthorized, "无token")
			c.Abort()
			return
		}

		log.Println("token: ", token)
		if !toolUtils.VerifyToken(token) {
			log.Println("无效token")
			sendHTTPError(c, http.StatusUnauthorized, "无效token")
			c.Abort()
			return
		}
		v, err := toolUtils.DecodeToken(token)
		if err != nil {
			log.Println(err)
			sendServerInternelError("解析token出错", c)
			c.Abort()
			return
		}
		var user = toolUtils.AuthUser{}
		err = json.Unmarshal([]byte(v), &user)
		if err != nil {
			log.Println(err)
			sendServerInternelError("解析token出错", c)
			c.Abort()
			return
		}
		c.Set(variable.AuthKey, user.UserName)
		c.Set(variable.PlatformKey, user.Platform)
	}
}
