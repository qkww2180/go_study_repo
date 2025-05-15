package jwt

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	header := DefautHeader
	payload := JwtPayload{
		ID:          "rj4t49tu49",
		Issue:       "微信",
		Audience:    "王者荣耀",
		Subject:     "购买道具",
		IssueAt:     time.Now().Unix(),
		Expiration:  time.Now().Add(2 * time.Hour).Unix(),
		UserDefined: map[string]any{"name": strings.Repeat("大乔乔", 1)}, //JwtPayload以明文形式在网络上传输，不要包含敏感信息，或者采用https对整个应用层报文加密
	}

	if token, err := GenJWT(header, payload); err != nil {
		fmt.Printf("生成json web token失败: %v", err)
	} else {
		fmt.Println(token)
		if _, p, err := VerifyJwt(token); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("JWT验证通过。欢迎 %s !\n", p.UserDefined["name"])
		}
	}
}

// go test -v .\f_encryption\jwt\ -run=^TestJWT$ -count=1
