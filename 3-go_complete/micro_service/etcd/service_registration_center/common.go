package service_registration_center

import (
	"errors"
	"net"
)

var (
	ETCD_CLUSTER = []string{"127.0.0.1:2379"}
)

//  /service/i_grpc/hello_service/ip1
//  /service/i_grpc/hello_service/ip2
//  /service/i_grpc/hello_service/ip3

//  /service/i_grpc/login_service/ip1
//  /service/i_grpc/login_service/ip2
//  /service/i_grpc/login_service/ip3

const (
	SERVICE_ROOT_PATH = "/service/i_grpc" //etcd key的前缀
	HELLO_SERVICE     = "hello_service"
)

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
