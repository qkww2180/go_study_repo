package f_encryption

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
)

// DesEncrypt DES加密
// 密钥必须是64位，所以key必须是长度为8的byte数组
func DesEncrypt(text string, key [8]byte) (string, error) {
	src := []byte(text)
	block, err := des.NewCipher(key[:]) //用des创建一个加密器cipher
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()      //分组的大小，blockSize=8
	src = PKCS7.Padding(src, blockSize) //填充

	//src  24   [0:7]   [8:24]   [8:15]
	//dst  24	[0:7]	[8:]
	out := make([]byte, len(src)) //密文和明文的长度一致
	dst := out
	for len(src) > 0 {
		//分组加密
		block.Encrypt(dst, src[:blockSize]) //对src进行加密，加密结果放到dst里
		//移到下一组
		src = src[blockSize:]
		dst = dst[blockSize:]
	}
	return hex.EncodeToString(out), nil
}

// DesDecrypt DES解密
// 密钥必须是64位，所以key必须是长度为8的byte数组
func DesDecrypt(text string, key [8]byte) (string, error) {
	src, err := hex.DecodeString(text) //转成[]byte
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		//分组解密
		block.Decrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	}
	out, _ = PKCS7.Unpadding(out, blockSize) //反填充
	return string(out), nil
}

func DesEncryptCBC(text string, key [8]byte) (string, error) {
	src := []byte(text)
	block, err := des.NewCipher(key[:]) //用des创建一个加密器cipher
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()      //分组的大小，blockSize=8
	src = PKCS7.Padding(src, blockSize) //填充

	out := make([]byte, len(src))                      //密文和明文的长度一致
	encrypter := cipher.NewCBCEncrypter(block, key[:]) //CBC分组模式加密
	encrypter.CryptBlocks(out, src)
	return hex.EncodeToString(out), nil
}

func DesDecryptCBC(text string, key [8]byte) (string, error) {
	src, err := hex.DecodeString(text) //转成[]byte
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	out := make([]byte, len(src))                      //密文和明文的长度一致
	encrypter := cipher.NewCBCDecrypter(block, key[:]) //CBC分组模式解密
	encrypter.CryptBlocks(out, src)
	out, _ = PKCS7.Unpadding(out, block.BlockSize()) //反填充
	return string(out), nil
}
