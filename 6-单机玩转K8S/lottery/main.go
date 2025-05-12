package main

import (
	"net/http"
	"os"

	"errors"
	"net"

	"github.com/gin-gonic/gin"
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
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet {
			if !ipNet.IP.IsLoopback() {
				if ipNet.IP.IsPrivate() { //内网地址
					// 跳过IPV6
					if ipNet.IP.To4() != nil {
						ipv4 = ipNet.IP.String()
						return
					}
				}
			}
		}
	}

	err = errors.New("ERR_NO_LOCAL_IP_FOUND")
	return
}

func main() {
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")
	engine.GET("/", func(ctx *gin.Context) {
		name := ctx.Query("name")       //从请求参数中取值
		project := os.Getenv("project") //从环境变量中取值
		ip, _ := GetLocalIP()
		ctx.HTML(http.StatusOK, "home.html", gin.H{"name": name, "project": project, "ip": ip})
	})
	port := os.Getenv("port") //从环境变量中取值
	engine.Run("0.0.0.0:" + port)
}
