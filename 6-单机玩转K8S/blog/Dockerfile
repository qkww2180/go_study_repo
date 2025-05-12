# 使用官方的 Golang 镜像作为基础镜像
FROM golang:latest

# 定义变量
ARG port=5678
ARG project=blog

# 定义环境变量
ENV project=$project
ENV port=$port

# 在容器内创建一个目录来存放我们的应用代码
RUN mkdir /app

# 将工作目录切换到 /app
WORKDIR /app

# 将当前目录下的所有文件拷贝到 /app 目录
COPY . .

# 编译 Go 应用程序
RUN go env -w GOPROXY=https://goproxy.io,direct   # 确保下载第三方库时走的是这个代理
RUN go build -o myapp .

# 暴露端口
EXPOSE $port

# 运行应用程序
CMD ["./myapp"]

