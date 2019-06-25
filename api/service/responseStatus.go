package service


type ApiStatus = int32
var (
	// common success
	APISuccess ApiStatus = 0
	// common fail
	APIFailed ApiStatus = 1
	// 请求参数错误
	APIParamsError ApiStatus = 2
	// 服务端错误
	APIServerInternalError ApiStatus = 3
)

