测试代码规范  
- 单元测试和基准测试必须放在以_test.go为后缀的文件里。
- 单元测试函数以Test开头，基准测试函数以Benchmark开头。
- 单元测试以*testing.T为参数，函数无返回值。
- 基准测试以*testing.B为参数，函数无返回值。  
# 单元测试
```shell
 go test -v ./b_project_prepare/test -run=^TestCopySlice$ -count=1 -timeout=20m
```
- -run 单测函数需要满足的正则表达式
- -v 打印详情测试信息
- -timeout 默认10分钟超时
- -count 函数运行几次  
- 把count设置1可以禁用缓存，否则主函数修改之后test结果可能不变

## FAIL和PASS
默认情况下单测都是PASS，除非以下几种情况： 
1. `t.Fail()`
2. `t.Error(args ...any)`  单测不会终止
3. `t.Fatal(args ...any)`  单测会终止
4. `assert`失败。想使用assert得借助于第三方库，比如`github.com/stretchr/testify/assert`

# 基础测试
```shell
 go test ./b_project_prepare/test -run=^$ -bench=CopySlice$ -benchmem -benchtime=3s -count=1 -cpuprofile=z_data/cpu -memprofile=z_data/mem
```
- -bench 基础测试函数需要满足的正则表达式
- -benchmem 输出内存开销的相关信息
- -benchtime 测试函数支行多长时间，默认是2秒
- -cpuprofile CPU pprof文件
- -memprofile 内存 pprof文件
# pprof
proof是可视化性能分析工具，提供以下功能：
1. CPU Profiling：按一定频率采集CPU使用情况。
2. Memory Profiling：监控内存使用情况，检查内存泄漏。
3. Goroutine Profiling：对正在运行的Goroutine进行堆栈跟踪和分析，检查协程泄漏。 

查看CPU使用情况`go tool pprof data/cpu`
进入交互界面后常用的命令有：  
- topn：列出最耗计算资源的前n个函数
- list func：列出某函数里每一行代码消耗多少计算资源
- peek func：列出某函数里最耗计算资源的前几个子函数   
- exit：退出交互界面

查看内存使用情况`go tool pprof data/mem`  
pprof结果可以在浏览器上进行可视化`go tool pprof -http=:8080 data/cpu`或`go tool pprof -http=:8080 data/mem`

# 测试覆盖率
`go test -cover $dir`   只能给出$dir目录的整体单测覆盖率

`go test ./project_prepare -coverprofile=data/test_cover`
或
`go test ./project_prepare -coverprofile=data/test_cover -covermode=count`
covermode的3个取值：
- set: 每个语句是否执行，默认值
- count: 每个语句执行了几次，鼠标悬停在语句上能显示执行的次数
- atomic: 类似于count, 但表示的是并行程序中的精确计数

`go tool cover -func=data/test_cover`		输出每一个函数的覆盖率

`go tool cover -html=data/test_cover`  	
细化到每一行代码的覆盖情况

