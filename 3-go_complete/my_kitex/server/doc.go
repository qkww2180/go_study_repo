package main

/*
kitex通过-service参数来生成服务端的脚手架代码。如果同名文件已存在，则不会被覆盖

cd server
kitex -module my_kitex -service dqq.math -use my_kitex/kitex_gen -I="D:\\go_project\\my_kitex" idl/math.proto
-service指定服务的名称，也是go build生成的可执行文件的名称
-use 参数表示让kitex不生成kitex_gen目录（因为这个目录之前已经生成过了），而使用该选项给出的 import path
-I指定proto文件的路径，windows上需要使用绝对路径，不能直接写../idl/math.proto
*/

/*
若报错 undefined: descriptorpb.Default_FileOptions_PhpGenericServices
是因为go.mod里的google.golang.org/protobuf v1.33.0和github.com/golang/protobuf v1.5.2不兼容，手动修改为github.com/golang/protobuf v1.5.4，或在go.mod末尾加一行 replace github.com/golang/protobuf => github.com/golang/protobuf v1.5.4

修改go.mod后要重新执行go mod tidy

若报错 not enough arguments in call to t.tProt.WriteMessageBegin
是因为 thrift 官方在 0.14 版本对 thrift 接口做了 breaking change，导致生成代码不兼容。在go.mod末尾加一行 replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

github.com/choleraehyq/pid和sonic每次go版本升级时都要更新到最新版本
*/

/*
完善handler.go
编译 go build -o output/bin/dqq.math.exe
*/
