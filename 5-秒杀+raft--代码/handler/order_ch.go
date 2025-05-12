package handler

import (
	"dqq/concurrency/database"
	"dqq/concurrency/util"
)

var (
	orderCh = make(chan database.Order, 10000)
	stopCh  = make(chan struct{}, 1)
)

func InitChannel() {
	go func() {
		<-stopCh
		close(orderCh)
	}()
}

// 把订单放入channel
func PutOrder(UserId, GiftId int) {
	order := database.Order{UserId: UserId, GiftId: GiftId}
	orderCh <- order
}

// 从channel里取出订单，把订单写入Mysql
func TakeOrder() {
	for {
		order, ok := <-orderCh
		if !ok {
			util.LogRus.Info("order channel is empty and closed")
			break
		}
		database.CreateOrder(order.UserId, order.GiftId) //写入mysql
	}
}

// 目的是想关闭orderCh，该函数可以反复调用，除第一次外，其他调用都要阻塞
func CloseChannel() {
	// stopCh <- struct{}{}

	select {
	case stopCh <- struct{}{}: //为了不让函数阻塞在本行代码，外面套一个select
		// util.LogRus.Info("close channel")
	default:
	}
}
