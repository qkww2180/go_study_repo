package util

import (
	"time"
)

var (
	loc, _ = time.LoadLocation("Asia/Shanghai")
)

// 即使上一次定时任务没执行完，下一次定时任务也会如期开始
func scheduleAsyn(work func(), hour int, minute int, second int) error {
	const period int64 = 86400
	ticker := time.NewTicker(1) //time.After，time.Ticker，time.Timer，time.Sleep都可以互相替换
	<-ticker.C
	defer ticker.Stop()

	now := time.Now()
	point := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, loc)
	diff := (point.Unix()%period + period - now.Unix()%period) % period
	// log.Printf("diff=%d\n", diff)
	if diff == 0 {
		go work()
	} else {
		ticker.Reset(time.Duration(diff) * time.Second)
		<-ticker.C
		go work()
	}

	ticker.Reset(time.Duration(period) * time.Second)
	for {
		<-ticker.C
		go work()
	}
}

// 如果上一个定时任务还没执行完，则跳过本次的定时任务
func scheduleSyn(work func(), hour int, minute int, second int) error {
	const period int64 = 86400
	ticker := time.NewTicker(1)
	<-ticker.C
	defer ticker.Stop()

	for {
		now := time.Now()
		point := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, loc)
		diff := (point.Unix()%period + period - now.Unix()%period) % period
		// log.Printf("diff=%d\n", diff)
		if diff == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		ticker.Reset(time.Duration(diff) * time.Second)
		<-ticker.C
		work()
	}
}

func Shedule(work func(), hour, minute, second int, asyn bool) {
	if asyn {
		scheduleAsyn(work, hour, minute, second)
	} else {
		scheduleSyn(work, hour, minute, second)
	}
}
