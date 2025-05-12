package algorithm

import "time"

// 基于令牌桶算法的限流器。不支持并发
type Limiter struct {
	rate   int64     //每秒钟生成令牌的速度
	cap    int       //桶的容量
	tokens float64   //目前还剩多少令牌，浮点数，而且可以为负
	last   time.Time //最后一次索取令牌的时间
}

func NewLimiter(rate int64, cap int) *Limiter {
	return &Limiter{
		rate:   rate,
		cap:    cap,
		tokens: 0,
		last:   time.Now(),
	}
}

// 生产这么多令牌需要多长时间
func (limiter *Limiter) durationFromTokens(tokens float64) time.Duration {
	n := tokens / float64(limiter.rate) //需要多少秒
	return time.Duration(float64(time.Second) * n)
}

func (limiter *Limiter) WaitN(n int) {
	now := time.Now()
	//从上一次索取令牌到现在，一共产生了多少新令牌
	delta := now.Sub(limiter.last).Seconds() * float64(limiter.rate)

	//当前最新的令牌数
	tokens := limiter.tokens + delta
	if tokens > float64(limiter.cap) {
		tokens = float64(limiter.cap)
	}

	// 减去要索取的令牌数
	tokens -= float64(n)
	limiter.tokens = tokens
	limiter.last = now //及时更新last字段
	if tokens < 0 {    //还缺少这么多令牌
		time.Sleep(limiter.durationFromTokens(-tokens)) //为了生成这么多令牌需要等这么长时间
	}
}
