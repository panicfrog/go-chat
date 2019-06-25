package toolUtils

import (
	"encoding/json"
	"fmt"
	"gochat/db"
	"gochat/variable"
	"time"
)

type AuthUser struct {
	db.User
	Platform variable.Platform `json:"platform"`
}

func (this *AuthUser) CreateToken() (token string, err error) {
	v, err := json.Marshal(this)
	if err != nil {
		return
	}
	var m = make(map[string]interface{})
	err = json.Unmarshal([]byte(v), &m)
	if err != nil {
		return
	}
	m["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano())
	b, err := json.Marshal(m)
	if err != nil {
		return
	}
	token = string(b)
	return
}