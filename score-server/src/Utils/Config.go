package Utils

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/gorilla/sessions"
)

var Sessions = sessions.NewCookieStore([]byte("testingkey"))

func HashMD5(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
