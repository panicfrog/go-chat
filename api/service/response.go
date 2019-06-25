package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func sendSuccess(message string, v interface{}, c *gin.Context) {
	sendResponse(c, APISuccess, message, v)
}

func sendFail(message string, c *gin.Context) {
	sendResponse(c, APIFailed, message, nil)
}

func sendParamError(message string, c *gin.Context) {
	sendResponse(c, APIParamsError, message, nil)
}

func sendServerInternelError(message string, c *gin.Context) {
	sendResponse(c, APIServerInternalError, message, nil)
}

func sendHTTPError(c *gin.Context,  state int, message string) {
	c.String(state, "%s", message)
}

func sendResponse(c *gin.Context, businessStateCode ApiStatus, message string, v interface{}) {
	res := ResponseStruct{businessStateCode, message, v}
	response, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	var m map[string]interface{}

	err = json.Unmarshal(response, &m)

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, m)
}
