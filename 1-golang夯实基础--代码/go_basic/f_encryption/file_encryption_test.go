package f_encryption_test

import (
	"dqq/go/basic/f_encryption"
	"fmt"
	"testing"
)

func TestFileEncryption(t *testing.T) {
	keyAES := []byte("ir489u58ir489u54") //AES算法key必须是长度为16的byte数组(128bit)。对称加密，加密的解密使用相同的key
	plainFile := "../z_data/verse.txt"

	encryptFileAES := "../z_data/verse.aes"
	plainFileAES := "../z_data/verse(解密aes).txt"
	if err := f_encryption.FileEncryption(plainFile, encryptFileAES, f_encryption.AES, keyAES); err != nil {
		fmt.Println(err)
	} else {
		if err = f_encryption.FileDecryption(encryptFileAES, plainFileAES, f_encryption.AES, keyAES); err != nil {
			fmt.Println(err)
		}
	}

	keyDES := []byte("ir489u58") //DES算法key必须是长度为8的byte数组(64bit)
	encryptFileDES := "../z_data/verse.des"
	plainFileDES := "../z_data/verse(解密des).txt"
	if err := f_encryption.FileEncryption(plainFile, encryptFileDES, f_encryption.DES, keyDES); err != nil {
		fmt.Println(err)
	} else {
		if err = f_encryption.FileDecryption(encryptFileDES, plainFileDES, f_encryption.DES, keyDES); err != nil {
			fmt.Println(err)
		}
	}
}

// go test -v ./f_encryption -run=^TestFileEncryption$ -count=1
