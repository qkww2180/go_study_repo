package f_encryption

import (
	"crypto/md5"
	"crypto/sha1"
)

func Sha1(data string) []byte {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return sha1.Sum(nil)
}

func Md5(data string) []byte {
	md5 := md5.New()
	md5.Write([]byte(data))
	return md5.Sum(nil)
}
