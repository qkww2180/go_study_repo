package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	enc_util "dqq/encryption/util"
	"dqq/util"
	"fmt"
	"io"
	"os"
)

const (
	SYMMETRICAL_ENCRYPTION = iota
	DES
	AES
)

// 文件加密
func FileEncryption(infile string, outfile string, algo int, key []byte) error {
	fin, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer fin.Close()
	fout, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fout.Close()

	content, err := io.ReadAll(fin)
	if err != nil {
		return err
	}

	var block cipher.Block
	switch algo {
	case AES:
		block, err = aes.NewCipher(key)
	case DES:
		block, err = des.NewCipher(key)
	default:
		return fmt.Errorf("unsurported encrypt algo %d", algo)
	}
	if err != nil {
		return err
	}

	blockSize := block.BlockSize()
	src := enc_util.PKCS5.Padding(content, blockSize) //加密算法的输入必须是blockSize的整倍数
	dest := make([]byte, len(src))                    //密文跟明文长度相同
	encrypter := cipher.NewCBCEncrypter(block, key)   //CBC分组模式加密
	encrypter.CryptBlocks(dest, src)                  //加密
	fout.Write(dest)
	return nil
}

// 文件解密
func FileDecryption(infile string, outfile string, algo int, key []byte) error {
	fin, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer fin.Close()
	fout, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fout.Close()

	content, err := io.ReadAll(fin)
	if err != nil {
		return err
	}

	var block cipher.Block
	switch algo {
	case AES:
		block, err = aes.NewCipher(key)
	case DES:
		block, err = des.NewCipher(key)
	default:
		return fmt.Errorf("unsurported encrypt algo %d", algo)
	}
	if err != nil {
		return err
	}

	decrypted := make([]byte, len(content))         //密文跟明文长度相同
	decrypter := cipher.NewCBCDecrypter(block, key) //CBC分组模式解密
	decrypter.CryptBlocks(decrypted, content)       //解密
	out, err := enc_util.PKCS5.Unpadding(decrypted, block.BlockSize())
	if err != nil {
		return err
	}
	fout.Write(out)
	return nil
}

func main() {
	keyAES := []byte("ir489u58ir489u54") //AES算法key必须是长度为16的byte数组(128bit)。对称加密，加密的解密使用相同的key
	plainFile := util.RootPath + "/data/学生信息表.xlsx"

	encryptFileAES := util.RootPath + "/data/学生信息表.aes"
	plainFileAES := util.RootPath + "/data/学生信息表(解密aes).xlsx"
	if err := FileEncryption(plainFile, encryptFileAES, AES, keyAES); err != nil {
		fmt.Println(err)
	} else {
		if err = FileDecryption(encryptFileAES, plainFileAES, AES, keyAES); err != nil {
			fmt.Println(err)
		}
	}

	keyDES := []byte("ir489u58") //DES算法key必须是长度为8的byte数组(64bit)
	encryptFileDES := util.RootPath + "/data/学生信息表.des"
	plainFileDES := util.RootPath + "/data/学生信息表(解密des).xlsx"
	if err := FileEncryption(plainFile, encryptFileDES, DES, keyDES); err != nil {
		fmt.Println(err)
	} else {
		if err = FileDecryption(encryptFileDES, plainFileDES, DES, keyDES); err != nil {
			fmt.Println(err)
		}
	}
}

// go run .\encryption\file\
