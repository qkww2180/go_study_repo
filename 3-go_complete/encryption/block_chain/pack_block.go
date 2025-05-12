package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

// 哈希算法
func Sha256(data string) []byte {
	sha := sha256.New()
	sha.Write([]byte(data))
	return sha.Sum([]byte(""))
}

// 区块
type Block struct {
	Index        int    //区块序号
	Timestamp    int64  //打包时间戳
	Random       int    //随机数
	Digest       []byte //本区块的哈希值
	PrevDigest   []byte //上一个区块的哈希值
	TradeRecords string //本区块里存储的所有交易记录
}

// 打包区块
func (self *Block) PackBlock(PrevBlock *Block) {
	self.Index = PrevBlock.Index + 1
	self.Timestamp = time.Now().Unix()
	self.PrevDigest = PrevBlock.Digest

	digest1 := Sha256(strconv.Itoa(self.Index) +
		strconv.FormatUint(uint64(self.Timestamp), 10) +
		string(self.PrevDigest) + self.TradeRecords) //先对固定部分（除随机数以外的部分）做哈希
	for {
		digest2 := Sha256(string(digest1) + strconv.Itoa(self.Random)) //对  固定部分+随机数  做哈希
		if digest2[0] == 0 {                                           //简单起见，这里只要求前8个bit是0
			self.Digest = digest2
			fmt.Printf("尝试了%d次\n", self.Random)
			break
		}
		self.Random += 1 //盲试，每次随机数加1
	}
}

// 验证区块
func VerifyBlock(block *Block) bool {
	digest1 := Sha256(strconv.Itoa(block.Index) + strconv.FormatUint(uint64(block.Timestamp), 10) + string(block.PrevDigest) + block.TradeRecords)
	digest2 := Sha256(string(digest1) + strconv.Itoa(block.Random))
	return slices.Equal(digest2, block.Digest)
}

func main() {
	trade := "A Transfer 10BTC to B"
	PrevBlock := Block{
		Index:  100,
		Digest: []byte("w8tu4t9i430q5="),
	}
	block := Block{
		TradeRecords: strings.Join([]string{trade, trade, trade}, ";"),
	}
	block.PackBlock(&PrevBlock)
	fmt.Println("当前区块的哈希值", hex.EncodeToString(block.Digest))
	fmt.Println("验证工作量", VerifyBlock(&block))
}

// go run .\encryption\block_chain\
