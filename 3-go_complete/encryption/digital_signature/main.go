package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

/**
生成1024位的RSA私钥：
openssl genrsa -out z_data/rsa_private_key.pem 1024
根据私钥生成公钥：
openssl rsa -in z_data/rsa_private_key.pem -pubout -out z_data/rsa_public_key.pem

pem是一种标准格式，它通常包含页眉和页脚
*/

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// go标准库不支持私钥加密，但直接提供了签名函数
func DigitalSignature(plain string) []byte {
	sha1 := sha1.New()
	sha1.Write([]byte(plain))
	digest := sha1.Sum([]byte("")) //参数一般置空即可，可以使用nil
	// 从文件中读取私钥
	privateKey, err := os.ReadFile("z_data/rsa_private_key.pem")
	checkError(err)
	block, _ := pem.Decode(privateKey)                           //解析PEM文件
	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes) //解析私钥。目前的数字证书一般都是基于ITU（国际电信联盟）制定的X.509标准
	checkError(err)
	//用私钥生成签名
	priv := privInterface.(*rsa.PrivateKey)
	signature, err := rsa.SignPKCS1v15(nil, priv, crypto.Hash(0), digest) //第一个和第三个参数一般置空即可
	checkError(err)
	return signature
}

// 验证数字签名
func VerifySignature(plain string, signature []byte) bool {
	sha1 := sha1.New()
	sha1.Write([]byte(plain))
	digest := sha1.Sum([]byte(""))
	// 从文件中读取公钥
	publicKey, err := os.ReadFile("z_data/rsa_public_key.pem")
	checkError(err)
	block, _ := pem.Decode(publicKey)                         //解析PEM文件
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes) //解析公钥
	checkError(err)
	pub := pubInterface.(*rsa.PublicKey)

	//用公钥验证签名
	return rsa.VerifyPKCS1v15(pub, crypto.Hash(0), digest, signature) == nil
}

func main() {
	plain := "因为我们没有什么不同"
	signature := DigitalSignature(plain)
	fmt.Println("验证数字签名", VerifySignature(plain, signature))
}

// go run .\f_encryption\digital_signature\
