package toolUtils

var key = "yeyongping123456"

func EncodeToken(v string) (token string, e error) {
	d, e := AESEncrypt([]byte(key), v)
	token = string(d)
	return
}

func DecodeToken(token string) (v string, e error) {
	d , e := AESDecrypt([]byte(key), token)
	v = string(d)
	return
}

func VerifyToken(token string) (verified bool) {
	verified = false
	_, e := DecodeToken(token)
	if e == nil {
		verified = true
	}
	return
}