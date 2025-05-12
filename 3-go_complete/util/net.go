package util

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"time"
)

// 判断一个ip是否为内网ip
func IsLocalNetIP(ipv4 string) bool {
	arr := strings.Split(ipv4, ".")
	if len(arr) != 4 {
		return false
	}
	brr := make([]int, 4)
	for i, ele := range arr {
		if n, err := strconv.Atoi(ele); err != nil {
			return false
		} else {
			brr[i] = n
		}
	}

	return brr[0] == 10 || // A类地址：10.0.0.0--10.255.255.255
		(brr[0] == 172 && brr[1] >= 16 && brr[1] <= 31) || // B类地址：172.16.0.0--172.31.255.255
		(brr[0] == 192 && brr[1] == 168) // C类地址：192.168.0.0--192.168.255.255
}

// ip转uint32
func Ip2Int(ipv4 string) uint32 {
	var rect uint32
	arr := strings.Split(ipv4, ".")
	if len(arr) != 4 {
		return 0
	}
	for i, ele := range arr {
		if n, err := strconv.Atoi(ele); err != nil {
			return 0
		} else {
			rect += (uint32(n) << uint32(8*(3-i)))
		}
	}
	return rect
}

// uint32转ip
func Int2Ip(n uint32) string {
	arr := make([]string, 4)
	var mask uint32 = 0xff
	for i := 0; i < 4; i++ {
		var shift uint32 = uint32(8 * (3 - i))
		ele := n & (mask << shift) >> shift
		arr[i] = strconv.Itoa(int(ele))
	}
	return strings.Join(arr, ".")
}

// 获取本机网卡IP(内网ip)
func GetLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}

	err = errors.New("ERR_NO_LOCAL_IP_FOUND")
	return
}

// 判断某个ip能否ping通
func Ping(ip string) bool {
	//ICMP即ping请求
	conn, err := net.DialTimeout("ip4:icmp", ip, time.Second*1)
	if err != nil {
		return false //连接不成功，有可能是对方设置了防火墙
	} else {
		conn.Close()
		return true //连接成功返回true
	}
}
