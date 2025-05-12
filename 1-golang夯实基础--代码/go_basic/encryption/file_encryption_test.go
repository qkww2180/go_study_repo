package encryption_test

import (
	"dqq/go/basic/encryption"
	"fmt"
	"testing"
)

func TestFileEncryption(t *testing.T) {
	keyAES := []byte("ir489u58ir489u54") //AES算法key必须是长度为16的byte数组(128bit)。对称加密，加密的解密使用相同的key
	plainFile := "../data/verse.txt"

	encryptFileAES := "../data/verse.aes"
	plainFileAES := "../data/verse(解密aes).txt"
	if err := encryption.FileEncryption(plainFile, encryptFileAES, encryption.AES, keyAES); err != nil {
		fmt.Println(err)
	} else {
		if err = encryption.FileDecryption(encryptFileAES, plainFileAES, encryption.AES, keyAES); err != nil {
			fmt.Println(err)
		}
	}

	keyDES := []byte("ir489u58") //DES算法key必须是长度为8的byte数组(64bit)
	encryptFileDES := "../data/verse.des"
	plainFileDES := "../data/verse(解密des).txt"
	if err := encryption.FileEncryption(plainFile, encryptFileDES, encryption.DES, keyDES); err != nil {
		fmt.Println(err)
	} else {
		if err = encryption.FileDecryption(encryptFileDES, plainFileDES, encryption.DES, keyDES); err != nil {
			fmt.Println(err)
		}
	}
}

// go test -v ./encryption -run=^TestFileEncryption$ -count=1
