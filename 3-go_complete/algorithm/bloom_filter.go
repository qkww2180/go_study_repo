package algorithm

import (
	"bufio"
	"encoding/gob"
	"log"
	"math/rand"
	"os"

	farm "github.com/dgryski/go-farm" //google FarmHash是CityHash的继承者
)

type BloomFilter struct {
	//成员变量必须可导出(exported)，否则无法序列到文件
	Arr      []byte   //用Byte数组表示底层bit数组
	BitCount uint     //几个bit
	Seeds    []uint32 //每个hash函数使用一个种子
}

// hashNum: hash函数的个数
//
// arrayLen: bit数组的长度
func NewBloomFilter(hashNum, arrayLen int) *BloomFilter {
	size := arrayLen / 8             //转为byte数组的长度
	seeds := make([]uint32, hashNum) //随机生成hashNum个种子，执行FarmHash时要用
	for i := 0; i < hashNum; i++ {
		seeds[i] = uint32(rand.Int())
	}
	return &BloomFilter{
		Arr:      make([]byte, size),
		BitCount: uint(arrayLen),
		Seeds:    seeds,
	}
}

// 判断bit数组的第index个元素是否为1
func (bf *BloomFilter) getBit(index uint) bool {
	var a uint = index / 8
	var bt uint = uint(bf.Arr[a])

	var b uint = index % 8
	var c uint = 1 << b
	return bt&c == c
}

// 把bit数组的第index个元素置为1
func (bf *BloomFilter) setBit(index uint) {
	var a uint = index / 8
	var b uint = index % 8
	var c uint = 1 << b
	bf.Arr[a] |= byte(c)
}

// 向BloomFilter中添加一个元素
func (bf *BloomFilter) Add(ele string) {
	bs := []byte(ele)
	for _, seed := range bf.Seeds {
		index := uint(farm.Hash32WithSeed(bs, seed)) % bf.BitCount
		bf.setBit(index)
	}
}

// 判断一个元素是否存在于BloomFilter中
func (bf *BloomFilter) Exists(ele string) bool {
	bs := []byte(ele)
	for _, seed := range bf.Seeds {
		index := uint(farm.Hash32WithSeed(bs, seed)) % bf.BitCount
		if !bf.getBit(index) {
			return false
		}
	}
	return true
}

// 把BloomFilter导出到文件
func (bf *BloomFilter) Dump(outfile string) error {
	fout, err := os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}
	defer fout.Close()
	writer := bufio.NewWriter(fout)
	defer writer.Flush()
	encoder := gob.NewEncoder(writer)
	return encoder.Encode(*bf)
}

// 从BloomFilter中加载文件
func LoadBloomFilter(infile string) *BloomFilter {
	fin, err := os.Open(infile)
	if err != nil {
		log.Printf("LoadBloomFilter from file %s failed: %v", infile, err)
		return nil
	}
	reader := bufio.NewReader(fin)
	decoder := gob.NewDecoder(reader)
	var bf BloomFilter
	if err := decoder.Decode(&bf); err != nil {
		log.Printf("DecodeBloomFilter from file %s failed: %v", infile, err)
		return nil
	} else {
		return &bf
	}
}
