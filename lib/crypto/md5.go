package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5(src string) string {
	h := md5.New()
	h.Write([]byte(src))
	data := hex.EncodeToString(h.Sum(nil))
	return data
}
