package f_encryption_test

import (
	"dqq/go/basic/f_encryption"
	"fmt"
	"testing"
)

func TestECC(t *testing.T) {
	prvKey, err := f_encryption.GenPrivateKey()
	if err != nil {
		t.Fatalf("genPrivateKey fail: %s\n", err)
	}
	pubKey := prvKey.PublicKey
	plain := "因为我们没有什么不同"
	cipher, err := f_encryption.ECCEncrypt(plain, pubKey)
	if err != nil {
		t.Fatalf("ECCEncrypt fail: %s\n", err)
	}
	plain, err = f_encryption.ECCDecrypt(cipher, prvKey)
	if err != nil {
		t.Fatalf("ECCDecrypt fail: %s\n", err)
	}
	fmt.Printf("明文：%s\n", plain)
}

// go test -v ./f_encryption -run=^TestECC$ -count=1
