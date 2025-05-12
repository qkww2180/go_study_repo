package test

import (
	"dqq/algorithm"
	"fmt"
	"testing"
	"time"
)

func TestTimeoutCache(t *testing.T) {
	tc := algorithm.NewTimeoutCache(10) //缓存容量为10

	tc.Add(1, "value", 1)   //1秒后到期
	tc.Add(2, "value", 100) //100秒后到期
	tc.Add(3, "value", 100)

	time.Sleep(2 * time.Second)

	for _, key := range []int{1, 2, 3} {
		_, exists := tc.Get(key)
		fmt.Printf("key %d exists %t\n", key, exists) //1到期不存在了，2 3还存在
	}
	fmt.Println()

	for i := 1; i <= 10; i++ { //再添加10个新元素
		tc.Add(i+3, "value", 200) //200秒后到期
	}

	for key := 1; key <= 13; key++ { //缓存最多只能容纳10个元素
		_, exists := tc.Get(key)
		fmt.Printf("key %d exists %t\n", key, exists) //2 3先到期不存在了，4-13还存在
	}
}

// go test ./algorithm/test -v -run=^TestTimeoutCache$ -count=1
