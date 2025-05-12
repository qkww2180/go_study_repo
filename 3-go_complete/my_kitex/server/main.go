package main

import (
	math_service "my_kitex/kitex_gen/math_service/math"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	// "github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
	klog.SetLevel(klog.LevelDebug)

	reg, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"}) //etcd需要用户名密码时，使用NewEtcdRegistryWithAuth。确保etcd已经启动好
	if err != nil {
		klog.Fatal(err.Error())
	}

	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:5678") //server的工作端口
	svr := math_service.NewServer(
		new(MathImpl),                //绑定service的具体实现
		server.WithServiceAddr(addr), //指定工作端口
		// server.WithACLRules(AuthMW),    //自定义访问控制
		server.WithMiddleware(TimerMW),     //可以添加多个中间件
		server.WithMiddleware(RecordArgMW), //添加中间件
		server.WithRegistry(reg),           //注册到etcd
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "dqq.math"}), //往etcd上注册时ServiceName不能为空
		// server.WithLimit(&limit.Option{MaxConnections: 100, MaxQPS: 600}),               //对grpc无效(对thrift有效)
	)
	// 查看是否注册成功 etcdctl get --prefix "kitex"
	// kill server之后，etcd上就没有相应的数据了

	err = svr.Run()

	if err != nil {
		klog.Error(err.Error())
	}
}

// go run ./server
