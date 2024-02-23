package md5

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func GetMd5(code string) string {
	//c1 := md5.Sum([]byte(code)) //返回[16]byte数组

	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}
