package main

import (
	"dqq/web_scraper/algorithm"
	"dqq/web_scraper/util"
	"log"
)

var (
	bloomFilter     *algorithm.BloomFilter
	bloomFilterFile string
)

func init() {
	bloomFilterFile = util.RootPath + "data/bili.bf"
}

func InitBloomFilter() {
	if ok, _ := util.PathExists(bloomFilterFile); ok {
		bloomFilter = algorithm.LoadBloomFilter(bloomFilterFile)
	}
	if bloomFilter == nil {
		bloomFilter = algorithm.NewBloomFilter(8, 4<<20)
	}
}

func DumpBloomFilter() {
	if bloomFilter != nil {
		if err := bloomFilter.Dump(bloomFilterFile); err != nil {
			log.Printf("dump job bloom filter to file %s failed: %v", bloomFilterFile, err)
		}
	}
}
