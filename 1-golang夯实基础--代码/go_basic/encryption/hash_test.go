package encryption_test

import (
	"dqq/go/basic/encryption"
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	data := "123456"
	hs := encryption.Sha1(data)
	fmt.Println("SHA-1", hs, len(hs))
	hm := encryption.Md5(data)
	fmt.Println("MD5", hm, len(hm))
}

// go test -v ./encryption -run=^TestHash$ -count=1
