package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	loc, _ = time.LoadLocation("Asia/Shanghai")
)

func parseInt(s string) (int, error) {
	if s == "*" {
		return 0, nil
	}
	return strconv.Atoi(s)
}

// cron形如：* 3 18，表示每小时的3分18秒执行定时任务
func parseCron(cron string) (int, int, int, int64, error) {
	arr := strings.Split(cron, " ")
	if len(arr) != 3 {
		return 0, 0, 0, 0, fmt.Errorf("invalid cron %s", cron)
	}
	hour, err := parseInt(arr[0])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid cron hour %s", arr[0])
	}
	minute, err := parseInt(arr[1])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid cron minute %s", arr[1])
	}
	second, err := parseInt(arr[2])
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid cron second %s", arr[2])
	}
	var period int64 = 86400
	if hour == 0 {
		period = 3600
		if minute == 0 {
			period = 60
			if second == 0 {
				return 0, 0, 0, 0, fmt.Errorf("invalid cron %s", cron)
			}
		}
	}
	return hour, minute, second, period, nil
}

// 即使上一次定时任务没执行完，下一次定时任务也会如期开始
func scheduleAsyn(work func(), cron string) error {
	hour, minute, second, period, err := parseCron(cron)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(1) //time.After，time.Ticker，time.Timer，time.Sleep都可以互相替换
	defer ticker.Stop()

	now := time.Now()
	point := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, loc)
	diff := (point.Unix()%period + period - now.Unix()%period) % period
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
func scheduleSyn(work func(), cron string) error {
	hour, minute, second, period, err := parseCron(cron)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(1)
	defer ticker.Stop()

	for {
		now := time.Now()
		point := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, loc)
		diff := (point.Unix()%period + period - now.Unix()%period) % period
		if diff == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		ticker.Reset(time.Duration(diff) * time.Second)
		<-ticker.C
		work()
	}
}

func Shedule(work func(), cron string, asyn bool) {
	if asyn {
		scheduleAsyn(work, cron)
	} else {
		scheduleSyn(work, cron)
	}
}

func main21() {
	go Shedule(func() {
		log.Println("schedule 1")
		time.Sleep(70 * time.Second)
	}, "* * 4", true)
	go Shedule(func() {
		log.Println("schedule 2")
		time.Sleep(70 * time.Second)

	}, "* * 10", false)
	time.Sleep(4 * time.Minute)
}
