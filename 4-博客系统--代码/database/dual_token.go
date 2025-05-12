package database

import (
	"blog/util"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

/**
什么时候使用Redis?
1. 高并发，低延时。redis比mysql快一到两个数量级。
2. redis可靠性没mysql高，万一redis挂了对业务影响不大，好修复。
3. redis通常存储string型value，此时它相对于mysql的性能优势更明显。
*/

const (
	TOKEN_PREFIX = "dual_token_"
	TOKEN_EXPIRE = 7 * 24 * time.Hour //一次登录7天有效
)

// 把<refreshToken, authToken>写入redis
func SetToken(refreshToken, authToken string) {
	client := GetRedisClient()
	if err := client.Set(TOKEN_PREFIX+refreshToken, authToken, TOKEN_EXPIRE).Err(); err != nil { //7天之后就拿不到authToken了
		util.LogRus.Errorf("write token pair(%s, %s) to redis failed: %s", refreshToken, authToken, err)
	}
}

// 根据refreshToken获取authToken
func GetToken(refreshToken string) (authToken string) {
	client := GetRedisClient()
	var err error
	if authToken, err = client.Get(TOKEN_PREFIX + refreshToken).Result(); err != nil {
		if err != redis.Nil {
			util.LogRus.Errorf("get auth token of refresh token %s failed: %s", refreshToken, err)
		}
	}
	return
}

// 该结构体映射为数据库里的一张表
type Art struct {
	Id   int
	Desc string
}

func GetArtById(id int) *Art {
	db := GetBlogDBConnection()
	var art Art
	if err := db.Select("id,desc").Where("id=?", id).First(&art).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			util.LogRus.Errorf("get art by id %d failed: %s", id, err)
		}
		return nil
	}
	//如果传id=0，art为空结构体并不能证明id=0的记录不存在
	return &art
}
