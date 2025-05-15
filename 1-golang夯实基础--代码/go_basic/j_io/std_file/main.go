package main

import (
	"fmt"
	"os"
)

// 从标准输入读入数据
func Scan() {
	fmt.Println("please input a word")
	var word string
	fmt.Scanf("%s\n", &word) //读入第1个空格前的单词，注意要加\n。不能走单元测试
	fmt.Printf("you input:[%s]\n", word)

	fmt.Println("please input two int")
	var word1 int
	var word2 int
	//整数用%d
	fmt.Scanf("%d %d\n", &word1, &word2) //读入多个单词，空格分隔。如果输入了更多单词会被缓存起来，丢给下一次scan
	fmt.Println("sum is:", word1+word2)

	fmt.Println("please input a line")
	content := make([]byte, 100)
	n, err := os.Stdin.Read(content)
	if err == nil {
		fmt.Print("you input:", string(content[:n]))
	}
}

// 标准输出 和 标准错误输出
func Stdout() {
	fmt.Println(os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd())
	fmt.Println("print") //输出到Stdout
	os.Stdout.WriteString("abc")
	os.Stderr.WriteString("123\n")
	fmt.Fprintf(os.Stdout, "%s\n", "hello") //输出到Stdout
	fmt.Fprint(os.Stderr, "golang")         //输出到Stderr
	// 操作系统级别（用C语言会看到下面的结果，但Go语言不是）
	// Stdin、Stdout、Stderr的文件句柄分别是0、1、2
	// Stdout是行缓冲，即遇到换行符才会输出到终端，在此之前一直是缓存在内存里；Stderr无缓冲，立即输出。
}

func main() {
	// Scan()
	Stdout()
}
