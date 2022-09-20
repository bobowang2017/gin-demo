package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encrypt(s, salt string) string {
	b := []byte(s)
	h := md5.New()
	h.Write(b)
	byteSalt := []byte(salt)
	h.Write(byteSalt)
	return hex.EncodeToString(h.Sum(nil))
}
