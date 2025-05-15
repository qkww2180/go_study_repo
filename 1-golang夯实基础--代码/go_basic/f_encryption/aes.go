package f_encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

func AesEncrypt(text string, key [16]byte) (string, error) {
	src := []byte(text)
	src = PKCS7.Padding(src, aes.BlockSize) //填充,AES的分组大小为16
	block, err := aes.NewCipher(key[:])     //构造函数，创建加密器
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCBCEncrypter(block, key[:]) //CBC分组模式加密
	encrypted := make([]byte, len(src))                //给密文申请内存空间
	encrypter.CryptBlocks(encrypted, src)              //加密
	return hex.EncodeToString(encrypted), nil
}

func AesDecrypt(text string, key [16]byte) (string, error) {
	src, err := hex.DecodeString(text) //转为[]byte
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key[:]) //构造函数，创建加密器
	if err != nil {
		return "", err
	}
	edecrypter := cipher.NewCBCDecrypter(block, key[:])     //CBC分组模式解密
	decrypted := make([]byte, len(src))                     //给明文申请内存空间
	edecrypter.CryptBlocks(decrypted, src)                  //解密
	out, _ := PKCS7.Unpadding(decrypted, block.BlockSize()) //反填充
	return string(out), nil
}
