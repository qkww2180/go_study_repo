package f_encryption_test

import (
	"dqq/go/basic/f_encryption"
	"fmt"
	"testing"
)

func TestRSA(t *testing.T) {
	f_encryption.ReadRSAKey("../z_data/rsa_public_key.pem", "../z_data/rsa_private_key.pem")

	plain := "因为我们没有什么不同"
	cipher, err := f_encryption.RsaEncrypt([]byte(plain))
	if err != nil {
		fmt.Println(err)
	} else {
		// fmt.Printf("密文：%s\n", hex.EncodeToString(cipher))
		fmt.Printf("密文：%v\n", (cipher))
		bPlain, err := f_encryption.RsaDecrypt(cipher)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("明文：%s\n", string(bPlain))
		}
	}
}

// go test -v ./f_encryption -run=^TestRSA$ -count=1
