package concurrence

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func Timeout1() {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		fmt.Println("这里是业务函数")
		time.Sleep(100 * time.Millisecond)
		done <- struct{}{}
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel() //调用cancel，触发Done
	}()
	select {
	case <-done:
		fmt.Println("业务函数未超时")
	case <-ctx.Done(): //ctx.Done()是一个管道，调用了cancel()都会关闭这个管道，然后读操作就会立即返回
		err := ctx.Err()            //如果发生Done（管道被关闭），Err返回Done的原因，可能是被Cancel了，也可能是超时了
		fmt.Println("业务函数超时:", err) //context canceled
	}
}

func Timeout2() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50) //超时后会自动调用context的Deadline，Deadline会触发Done
	defer cancel()
	done := make(chan struct{})

	go func() {
		fmt.Println("这里是业务函数")
		time.Sleep(100 * time.Millisecond)
		done <- struct{}{}
	}()

	select {
	case <-done:
		fmt.Println("业务函数未超时")
	case <-ctx.Done(): //ctx.Done()是一个管道，context超时或者调用了cancel()都会关闭这个管道，然后读操作就会立即返回
		err := ctx.Err()            //如果发生Done（管道被关闭），Err返回Done的原因，可能是被Cancel了，也可能是超时了
		fmt.Println("业务函数超时:", err) //context deadline exceeded
	}
}

// 通过context.WithTimeout创建的Context，其寿命不会超过父Context的寿命。
func InheritTimeout() {
	parent, cancel1 := context.WithTimeout(context.Background(), time.Millisecond*1000) //parent设置100ms超时
	t0 := time.Now()
	defer cancel1()

	time.Sleep(500 * time.Millisecond) //消耗掉500ms

	// child, cancel2 := context.WithTimeout(parent, time.Millisecond*1000) //parent还剩500ms，child设置了1000ms之后到期，child.Done()管道的关闭时刻以较早的为准，即500ms后到期
	child, cancel2 := context.WithTimeout(parent, time.Millisecond*100) //parent还剩500ms，child设置了100ms之后到期，child.Done()管道的关闭时刻以较早的为准，即100ms后到期
	t1 := time.Now()
	defer cancel2()

	<-child.Done()
	fmt.Println(time.Since(t0).Milliseconds(), time.Since(t1).Milliseconds())
	fmt.Println(child.Err()) //context deadline exceeded
}

type StringKey string
type StringValue string
type IntValue int

func RoutineID() {
	for i := 0; i < 3; i++ {
		ctx := context.WithValue(context.Background(), StringKey("gid"), IntValue(i))
		ctx = context.WithValue(ctx, StringKey("owner"), StringValue("dqq"))
		go func(ctx context.Context) {
			if gid, ok := ctx.Value(StringKey("gid")).(IntValue); ok {
				fmt.Printf("本协程ID %d\n", gid)
			}
			if owner, ok := ctx.Value(StringKey("owner")).(StringValue); ok {
				fmt.Printf("owner %s\n", owner)
			}
			if name, ok := ctx.Value(StringKey("name")).(StringValue); ok {
				fmt.Printf("name %s\n", name)
			} else {
				fmt.Println("name不存在")
			}
		}(ctx)
	}
	time.Sleep(time.Second)
}

func step1(ctx *context.Context) {
	//根据父context创建子context，创建context时允许设置一个<key,value>对，key和value可以是任意数据类型
	*ctx = context.WithValue(*ctx, StringKey("name"), StringValue("大脸猫"))
}

func step2(ctx *context.Context) {
	//子context继承了父context里的所有key value
	*ctx = context.WithValue(*ctx, StringKey("age"), IntValue(18))
}

func step3(ctx context.Context) {
	if name, ok := ctx.Value(StringKey("name")).(StringValue); ok { //取出key对应的value，把any断言为string
		fmt.Printf("name %s\n", name)
	}

	if age, ok := ctx.Value(StringKey("age")).(IntValue); ok { //取出key对应的value，把any断言为int
		fmt.Printf("age %d\n", age)
	}
}

// 用context携带value仅用于跨进程传输数据和API调用，不用于向函数传递可选参数
func ContextWithValue() {
	ctx := context.Background() //空context
	step1(&ctx)                 //father里有一对<key,value>
	step2(&ctx)                 //grandson里有两对<key,value>
	step3(ctx)

	// 下面才是使用context携带value的正规场景
	client := http.Client{Timeout: 2 * time.Second}
	request, _ := http.NewRequest("GET", "http://127.0.0.1:1234", nil)
	client.Do(request) // request.Context()里会携带超时时间
}
