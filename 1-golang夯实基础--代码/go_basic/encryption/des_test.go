package encryption_test

import (
	"dqq/go/basic/encryption"
	"fmt"
	"log"
	"testing"
)

func TestDES(t *testing.T) {
	key := [8]byte{34, 65, 12, 125, 65, 70, 54, 27} //key必须是长度为8的byte数组
	plain := "因为我们没有什么不同"
	cipher, err := encryption.DesEncrypt(plain, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("密文：%s\n", cipher)

	plain, err = encryption.DesDecrypt(cipher, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("明文：%s\n", plain)
	fmt.Println("-------------------------------------")

	cipher, _ = encryption.DesEncryptCBC(plain, key)
	fmt.Printf("密文：%s\n", cipher)
	plain, _ = encryption.DesDecryptCBC(cipher, key)
	fmt.Printf("明文：%s\n", plain)
}

// go test -v ./encryption -run=^TestDES$ -count=1
