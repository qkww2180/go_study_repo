package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"io"
	"os"
)

const (
	_ = iota
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

	content, err := io.ReadAll(fin) //一次性读取文件里的所有内容
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
	encrypter := cipher.NewCBCEncrypter(block, key) //CBC分组模式加密

	src := PKCS7.Padding(content, block.BlockSize()) //加密算法的输入必须是blockSize的整倍数
	dest := make([]byte, len(src))                   //密文跟明文长度相同
	encrypter.CryptBlocks(dest, src)                 //加密
	fout.Write(dest)                                 //密文写入文件
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

	content, err := io.ReadAll(fin) //一次性读取文件里的所有内容
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
	decrypter := cipher.NewCBCDecrypter(block, key) //CBC分组模式解密

	decrypted := make([]byte, len(content))   //密文跟明文长度相同
	decrypter.CryptBlocks(decrypted, content) //解密
	out, err := PKCS7.Unpadding(decrypted, block.BlockSize())
	if err != nil {
		return err
	}
	fout.Write(out) //明文写入文件
	return nil
}
