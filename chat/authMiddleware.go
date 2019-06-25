package chat

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gochat/toolUtils"
	"gochat/variable"
	"log"
)

// TODO: 还有平台信息没有处理

func ChatAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.URL.Query().Get("token")
		if token == "" {
			c.Next()
			log.Println("未授权的websocket连接，reason: token不存在")
			return
		}
		//token = url.PathEscape(token)
		if !toolUtils.VerifyToken(token) {
			c.Next()
			log.Println("未授权的websocket连接，reason: '" + token + "'验证不通过")
			return
		}
		v, err := toolUtils.DecodeToken(token)
		if err != nil {
			c.Next()
			log.Println(err)
			return
		}
		var user = toolUtils.AuthUser{}
		err = json.Unmarshal([]byte(v), &user)
		if err != nil {
			c.Next()
			log.Println(err)
			return
		}
		c.Set(variable.AuthKey, user.UserName)
		c.Set(variable.PlatformKey, user.Platform)
	}
}