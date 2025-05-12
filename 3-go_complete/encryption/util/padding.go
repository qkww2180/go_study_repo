package util

import (
	"bytes"
	"errors"
)

var (
	ErrPaddingSize = errors.New("padding size error")
)

// pkcs7padding和pkcs5padding的填充方式相同，填充字节的值都等于填充字节的个数。例如需要填充4个字节，则填充的值为"4 4 4 4"。
var (
	PKCS5 = &pkcs5{}
)
var (
	// difference with pkcs5 only block must be 8
	PKCS7 = &pkcs5{}
)

// pkcs5Padding is a pkcs5 padding struct.
type pkcs5 struct{}

// Padding implements the Padding interface Padding method.
func (p *pkcs5) Padding(src []byte, blockSize int) []byte {
	srcLen := len(src)
	padLen := blockSize - (srcLen % blockSize) //注意： 当srcLen是blockSize的整倍数时，padLen等于blockSize而非0
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padText...)
}

// Unpadding implements the Padding interface Unpadding method.
func (p *pkcs5) Unpadding(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, ErrPaddingSize
	}
	return src[:srcLen-paddingLen], nil
}
