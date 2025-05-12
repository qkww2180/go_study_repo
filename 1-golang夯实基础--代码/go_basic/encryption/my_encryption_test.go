package encryption_test

import (
	"bytes"
	"dqq/go/basic/encryption"
	"fmt"
	"testing"
)

func TestMyEncryption(t *testing.T) {
	key := [8]byte{34, 65, 12, 125, 65, 70, 54, 27}

	algo := encryption.NewMyEncryption(key, encryption.NONE)
	plain := []byte("明月多情应笑我")
	cypher := algo.Encrypt(plain)
	fmt.Println(cypher)
	plain2, err := algo.Decrypt(cypher)
	if err != nil {
		t.Error(err)
	} else {
		if !bytes.Equal(plain, plain2) { // 比较两个byte切片里的元素是否完全相等
			fmt.Println(len(plain2), string(plain2))
			t.Fail()
		}
	}

	algo = encryption.NewMyEncryption(key, encryption.CBC)
	plain = []byte("明月多情应笑我")
	cypher = algo.Encrypt(plain)
	fmt.Println(cypher)
	plain2, err = algo.Decrypt(cypher)
	if err != nil {
		t.Error(err)
	} else {
		if !bytes.Equal(plain, plain2) { // 比较两个byte切片里的元素是否完全相等
			fmt.Println(len(plain2), string(plain2))
			t.Fail()
		}
	}
}

// go test -v ./encryption -run=^TestMyEncryption$ -count=1
