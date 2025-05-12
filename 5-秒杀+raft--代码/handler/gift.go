package handler

import (
	"dqq/concurrency/database"
	"dqq/concurrency/util"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取所有奖品信息，用于初始化轮盘
func GetAllGifts(ctx *gin.Context) {
	gifts := database.GetAllGiftsV1()
	if len(gifts) == 0 {
		ctx.JSON(http.StatusInternalServerError, nil)
	} else {
		//抹掉敏感信息
		for _, gift := range gifts {
			gift.Count = 0
		}
		ctx.JSON(http.StatusOK, gifts)
	}
}

// 抽奖
func Lottery(ctx *gin.Context) {
	for try := 0; try < 10; try++ { //最多重试10次
		gifts := database.GetAllGiftInventory() //获取所有奖品剩余的库存量
		ids := make([]int, 0, len(gifts))
		probs := make([]float64, 0, len(gifts))
		for _, gift := range gifts {
			if gift.Count > 0 { //先确保redis返回的库存量大小0，因为抽奖算法Lottery不支持抽中概率为0的奖品
				ids = append(ids, gift.Id)
				probs = append(probs, float64(gift.Count))
			}
		}
		if len(ids) == 0 {
			if os.Getenv("queue") == "channel" {
				CloseChannel() //关闭channel
			} else {
				go CloseMQ() //关闭写mq的连接
			}
			ctx.String(http.StatusOK, strconv.Itoa(0)) //0表示所有奖品已抽完
			return
		}
		index := util.Lottery(probs) //抽中第index个奖品
		giftId := ids[index]
		err := database.ReduceInventory(giftId) // 先从redis上减库存
		if err != nil {
			util.LogRus.Warnf("奖品%d减库存失败", giftId) //设想，某奖品只剩1件，并发情况下多个协程恰好都抽中了该奖品，第一个协程减库存后为0，第一个协程减库存后为负数--即减库存失败，即本次抽奖失败，进入下一轮for循环重试
			continue                               //减库存失败，则重试
		} else {
			if os.Getenv("queue") == "channel" {
				// 用户ID写死为1，关于用户身份认证参考《双Token认证博客系统》
				PutOrder(1, giftId) //把订单信息写入channel
			} else {
				ProduceOrder(1, giftId) //把订单信息写入mq
			}
			ctx.String(http.StatusOK, strconv.Itoa(giftId)) //减库存成功后才给前端返回奖品ID
			return
		}
	}
	ctx.String(http.StatusOK, strconv.Itoa(database.EMPTY_GIFT)) //如果10次之后还失败，则返回“谢谢参与”
}
