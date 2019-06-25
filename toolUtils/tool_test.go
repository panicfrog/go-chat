package toolUtils

import (
	"encoding/json"
	"gochat/db"
	"gochat/variable"
	"testing"
)

func TestAes(t *testing.T)  {
	var key = "1234567890134446"
	var a = "{\"passwd\":\"你是谁\",\"platform\":\"iOS\",\"timestamp\":\"1561102975452257000\",\"user_name\":\"woshishui\"}"
	v, e := AESEncrypt([]byte(key), a)
	if e != nil {
		t.Error(e)
	}
	m , e := AESDecrypt([]byte(key), v)
	if e != nil {
		t.Error(e)
	}
	if m != a {
		t.Error("不相等")
	}
}

func TestCreateToken(t *testing.T) {
	var user = AuthUser{db.User{UserName: "woshishui", Passwd:"你是谁"}, variable.PlatformIOS}
	token, err := user.CreateToken()
	if err != nil {
		t.Error(err)
		return
	}
	var _user = AuthUser{}
	err = json.Unmarshal([]byte(token), &_user)
	if err != nil {

		return
	}

	if _user.UserName != user.UserName {
		t.Error("json解析出错")
	}

	tk, ee := EncodeToken(token)
	if  ee != nil {
		t.Error(ee)
		return
	}

	_token, de := DecodeToken(tk)
	if de != nil {
		t.Error(de)
		return
	}

	if  _token != token {
		t.Error("验证出错")
	}

	if de := VerifyToken(tk); !de {
		t.Error("未通过")
		return
	}


}