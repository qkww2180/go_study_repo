package f_encryption_test

import (
	"dqq/go/basic/f_encryption"
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	data := "123456"
	hs := f_encryption.Sha1(data)
	fmt.Println("SHA-1", hs, len(hs))
	hm := f_encryption.Md5(data)
	fmt.Println("MD5", hm, len(hm))
}

// go test -v ./f_encryption -run=^TestHash$ -count=1
