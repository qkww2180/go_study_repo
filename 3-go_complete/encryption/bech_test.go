package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"dqq/encryption/util"
	"encoding/pem"
	"os"
	"testing"
)

var (
	content = []byte("At present, you need to have a FreeBSD, Linux, macOS, or Windows machine to run Go. We will use $ to represent the command prompt.")
	key     = []byte("ir489u58ir489u54")
)

func BenchmarkAES(b *testing.B) {
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encrypter := cipher.NewCBCEncrypter(block, key) //CBC分组模式加密
		src := util.PKCS5.Padding(content, blockSize)   //加密算法的输入必须是blockSize的整倍数
		dest := make([]byte, len(src))                  //密文跟明文长度相同
		encrypter.CryptBlocks(dest, src)                //加密
	}
}

func BenchmarkRSA(b *testing.B) {
	publicKey, _ := os.ReadFile("../data/rsa_public_key.pem") //相对于encryption的路径，因为go test指定的路径是encryption
	block, _ := pem.Decode(publicKey)
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	pub := pubInterface.(*rsa.PublicKey)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rsa.EncryptPKCS1v15(rand.Reader, pub, content)
	}
}

func BenchmarkHMACSHA256(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := hmac.New(sha256.New, key)
		h.Write(content)
		h.Sum(nil)
	}
}

// go test ./encryption -bench=BenchmarkAES -run=^$ -count=1 -timeout=2s -benchmem
// go test ./encryption -bench=BenchmarkRSA -run=^$ -count=1 -timeout=2s -benchmem
// go test ./encryption -bench=BenchmarkHMACSHA256 -run=^$ -count=1 -timeout=2s -benchmem
// go test ./encryption -bench=^Benchmark -run=^$ -count=1 -timeout=2s -benchmem

/**
BenchmarkAES-8                   2743167               388.7 ns/op           560 B/op          6 allocs/op
BenchmarkRSA-8                  13628943                87.81 ns/op            0 B/op          0 allocs/op
BenchmarkHMACSHA256-8            1837024               623.6 ns/op           512 B/op          6 allocs/op
**/
