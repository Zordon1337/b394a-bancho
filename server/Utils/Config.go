package Utils

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/gorilla/sessions"
)

var Sessions = sessions.NewCookieStore([]byte("testingkey"))

func CanRegister() bool {
	return false // here you can set if you want to allow registering
}
func HashMD5(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
