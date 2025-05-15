package algorithm

import (
	"dqq/web_scraper/util"
	"fmt"
	"testing"
)

func TestBloomFilter(t *testing.T) {
	bf := NewBloomFilter(8, 1<<20) //8个hash函数，1M个bit
	a, b, c, d := "野径人俱黑", "江船火独明", "晓看红湿处", "花重锦官城"
	bf.Add(a)
	bf.Add(b)
	if bf.Exists(c) {
		t.Fail()
	}
	if bf.Exists(d) {
		t.Fail()
	}
	if !bf.Exists(a) {
		t.Fail()
	}
	if !bf.Exists(b) {
		t.Fail()
	}

	//导出到文件，再从文件导入
	file := util.RootPath + "z_data/bloom_filter.bin"
	if err := bf.Dump(file); err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		bf2 := LoadBloomFilter(file)
		if bf2 == nil {
			t.Fail()
		} else {
			if bf2.Exists(c) {
				t.Fail()
			}
			if bf2.Exists(d) {
				t.Fail()
			}
			if !bf2.Exists(a) {
				t.Fail()
			}
			if !bf2.Exists(b) {
				t.Fail()
			}
		}
	}
}
