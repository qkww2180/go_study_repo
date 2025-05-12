package encryption_test

import (
	"dqq/go/basic/encryption"
	"fmt"
	"log"
	"testing"
)

func TestAES(t *testing.T) {
	key := [16]byte{'i', 'r', '4', '8', '9', 'u', '5', '8', 'i', 'r', '4', '8', '9', 'u', '5', '4'} //key必须是长度为16的byte数组
	plain := "因为我们没有什么不同"
	cipher, err := encryption.AesEncrypt(plain, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("密文：%s\n", cipher)

	plain, err = encryption.AesDecrypt(cipher, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("明文：%s\n", plain)
}

// go test -v ./encryption -run=^TestAES$ -count=1
