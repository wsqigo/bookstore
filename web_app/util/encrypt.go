package util

import (
	"crypto/md5"
	"encoding/hex"
)

const secret = "wsqigo"

func EncryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))

	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
