package encryption

import "errors"

const (
	NONE = iota //0
	CBC         //1
)

// XOR 异或运算，要求plain和key的长度相同
//
// 单看一个比特，任何数(0或1)跟0异或还是自己、跟1异或是其反面，所以任何数跟0两次异或还是自己、任何数跟1两次异或还是自己
type MyEncryption struct {
	key       [8]byte
	BlockSize int
	BlockMode int
}

func NewMyEncryption(key [8]byte, BlockMode int) *MyEncryption {
	return &MyEncryption{
		key:       key,
		BlockSize: 8,
		BlockMode: BlockMode,
	}
}

// 加密一组数据。plain和cypher的长度一定都等于BlockSize
func (en *MyEncryption) encryptBlock(plain, cypher []byte) {
	for j := 0; j < len(plain); j++ {
		cypher[j] = plain[j] ^ en.key[j] //不开辟额外的内存空间，就地加密
	}
}

// 解密一组数据。plain和cypher的长度一定都等于BlockSize
func (en *MyEncryption) decryptBlock(plain, cypher []byte) {
	for j := 0; j < len(plain); j++ {
		plain[j] = cypher[j] ^ en.key[j] //不开辟额外的内存空间，就地解密
	}
}

// 分组混淆。plain和prevCypher的长度一定相等
func (en *MyEncryption) confuseBlock(plain, prevCypher []byte) {
	switch en.BlockMode {
	case CBC:
		for j := 0; j < len(plain); j++ {
			plain[j] = plain[j] ^ prevCypher[j] //不开辟额外的内存空间，就地混淆
		}
	default:
	}
}

func (en *MyEncryption) deconfuseBlock(plain, prevCypher []byte) {
	en.confuseBlock(plain, prevCypher)
}

func (en *MyEncryption) Encrypt(plain []byte) []byte {
	plainPadding := PKCS7.Padding(plain, en.BlockSize) //明文末尾填充字节，长度成为BlockSize的整倍数
	cypher := make([]byte, len(plainPadding))
	prevCypher := make([]byte, en.BlockSize) //全0，任何数跟0异或还是自身
	for i := 0; i < len(plainPadding); i += en.BlockSize {
		begin := i
		end := i + en.BlockSize
		en.confuseBlock(plainPadding[begin:end], prevCypher)
		en.encryptBlock(plainPadding[begin:end], cypher[begin:end])
		copy(prevCypher, cypher[begin:end])
	}
	return cypher
}

func (en *MyEncryption) Decrypt(cypher []byte) ([]byte, error) {
	if len(cypher)%en.BlockSize != 0 {
		return nil, errors.New("密文长度不合法")
	}
	if len(cypher) == 0 {
		return []byte{}, nil
	}
	plainPadding := make([]byte, len(cypher))
	blockNum := len(cypher) / en.BlockSize //几组
	// 倒着解密，先解密最后一个分组
	for i := blockNum - 1; i >= 0; i-- {
		begin := i * en.BlockSize
		end := begin + en.BlockSize
		var prevCypher []byte //前一组密文
		if i == 0 {
			prevCypher = make([]byte, en.BlockSize) //全0
		} else {
			prevCypher = cypher[begin-en.BlockSize : end-en.BlockSize]
		}
		en.decryptBlock(plainPadding[begin:end], cypher[begin:end]) //解密得到明文，注意这还不是真正的明文
		en.deconfuseBlock(plainPadding[begin:end], prevCypher)      //跟上一组密文解混淆，得到的才是真正的明文
	}
	return PKCS7.Unpadding(plainPadding, en.BlockSize) //去除末尾填充的字节
}
