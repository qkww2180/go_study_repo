package redis_class

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
)

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", //没有密码
		DB:       0,  //redis默认会创建0-15号DB，这里使用默认的DB
	})
}

func TestStringValue(t *testing.T) {
	stringValue(context.Background(), client)
}

func TestListValue(t *testing.T) {
	listValue(context.Background(), client)
}

func TestSetgValue(t *testing.T) {
	setValue(context.Background(), client)
}

func TestZSetValue(t *testing.T) {
	zsetValue(context.Background(), client)
}

func TestHashTableValue(t *testing.T) {
	hashtableValue(context.Background(), client)
}

// go test -v ./database/redis -run=^TestStringValue$ -count=1
// go test -v ./database/redis -run=^TestListValue$ -count=1
// go test -v ./database/redis -run=^TestSetgValue$ -count=1
// go test -v ./database/redis -run=^TestZSetValue$ -count=1
// go test -v ./database/redis -run=^TestHashTableValue$ -count=1
