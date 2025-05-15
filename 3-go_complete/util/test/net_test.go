package test

import (
	"bufio"
	"dqq/util"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func TestIsLocalNetIP(t *testing.T) {
	if !util.IsLocalNetIP("172.16.0.235") {
		t.Fail()
	}
	if !util.IsLocalNetIP("172.16.0.10") {
		t.Fail()
	}
	if !util.IsLocalNetIP("172.16.0.193") {
		t.Fail()
	}
	if util.IsLocalNetIP("1.1.1.1") {
		t.Fail()
	}
}

func TestIp2Int(t *testing.T) {
	for i := 0; i < 10000; i++ {
		n1 := uint32(rand.Int63())
		ip := util.Int2Ip(n1)
		n2 := util.Ip2Int(ip)
		if n1 != n2 {
			fmt.Println(n1, ip, n2)
			t.Fail()
			break
		}
	}
}

func TestGenRandomIp(t *testing.T) {
	fout, err := os.OpenFile(util.RootPath+"z_data/ip_topk/ip.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}
	defer fout.Close()
	writer := bufio.NewWriter(fout)
	defer writer.Flush()
	for i := 0; i < 100000; i++ {
		n1 := uint32(rand.Int63n(1000) + 100000000000)
		ip := util.Int2Ip(n1)
		writer.WriteString(ip + "\n")
	}
}

func TestGetLocalIP(t *testing.T) {
	fmt.Println(util.GetLocalIP())
}

// go test -v .\util\test\ -run=TestIsLocalNetIP -count=1
// go test -v .\util\test\ -run=TestIp2Int -count=1
// go test -v .\util\test\ -run=TestGenRandomIp -count=1
// go test -v .\util\test\ -run=TestGetLocalIP -count=1
