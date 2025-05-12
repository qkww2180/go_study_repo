package encryption

import (
	ecies "github.com/ecies/go/v2"
)

func GenPrivateKey() (*ecies.PrivateKey, error) {
	return ecies.GenerateKey()
}

// ECCEncrypt 椭圆曲线加密
func ECCEncrypt(plain string, pubKey *ecies.PublicKey) ([]byte, error) {
	src := []byte(plain)
	return ecies.Encrypt(pubKey, src)
}

// ECCDecrypt 椭圆曲线解密
func ECCDecrypt(cipher []byte, prvKey *ecies.PrivateKey) (string, error) {
	if src, err := ecies.Decrypt(prvKey, cipher); err != nil {
		return "", err
	} else {
		return string(src), nil
	}
}
