package encryption

import (
	"bytes"
	"errors"
)

var (
	ErrPaddingSize = errors.New("padding size error")
)

var (
	PKCS7 = &pkcs7{}
)

type pkcs7 struct{}

// 填充字节的值都等于填充字节的个数。例如需要填充4个字节，则填充的值为"4 4 4 4"。PKCS5的blockSize固定为8
func (p *pkcs7) Padding(src []byte, blockSize int) []byte {
	srcLen := len(src)
	padLen := blockSize - (srcLen % blockSize) //注意： 当srcLen是blockSize的整倍数时，padLen等于blockSize而非0
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(src, padText...)
}

func (p *pkcs7) Unpadding(src []byte, blockSize int) ([]byte, error) {
	srcLen := len(src)
	paddingLen := int(src[srcLen-1])
	if paddingLen >= srcLen || paddingLen > blockSize {
		return nil, ErrPaddingSize
	}
	return src[:srcLen-paddingLen], nil
}
