package service

type ResponseStruct struct {
	Sc      ApiStatus   `json:"sc"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
