package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// UTF-8 转 GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 逐行读取文本文件
func readTextFileByLine(file string) {
	fin, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	reader := bufio.NewReader(fin)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			} else {
				if len(line) > 0 {
					fmt.Println(string(line))
				}
				break
			}

		} else {
			fmt.Println(string(line))
		}
	}
}

// 把文本文件分割成4段
func splitTextFile(file string) {
	fin, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	const count = 4.0
	info, err := fin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := info.Size()                          //获取文件大小
	seg := int(math.Ceil(float64(size) / count)) //分割成4段，计算每段大小
	fmt.Printf("文件总大小: %d，每段大小: %d\n", size, seg)
	//对于文本文件，没有存储额外的元信息，存储的全部都是内容数据本身。
	//同一个字符，不同的编码方式将其转为不同的byte数字。

	for {
		buffer := make([]byte, seg)
		n, err := fin.Read(buffer) //n是成功读取的字节数
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			} else {
				if len(buffer) > 0 {
					fmt.Println(buffer[:n])
					fmt.Print(string(buffer[:n]))
				}
				break
			}
		}
		fmt.Println(buffer[:n])
		fmt.Print(string(buffer[:n]))
	}
}

// 把二进制文件分割成4段
func splitBinaryFile(file string) []string {
	baseName := file
	suffix := ""
	pos := strings.LastIndex(file, ".")
	if pos >= 0 {
		baseName = file[:pos]
		suffix = file[pos:]
	}

	fin, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fin.Close()

	const count = 4.0
	info, err := fin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := info.Size()                          //获取文件大小
	seg := int(math.Ceil(float64(size) / count)) //分割成4段，计算每段大小
	fmt.Printf("文件总大小: %d，每段大小: %d\n", size, seg)
	//对于二进制文件，文件开头存储了一些元信息，比如文件格式、文件总大小等

	files := make([]string, 0, count)
	for i := 0; i < int(count); i++ {
		outFile := fmt.Sprintf("%s%d%s", baseName, i+1, suffix)
		files = append(files, outFile)
		fout, err := os.OpenFile(outFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		defer fout.Close()

		buffer := make([]byte, seg)
		n, err := fin.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			} else {
				if len(buffer) > 0 {
					fout.Write(buffer[:n])
				}
				break
			}
		}
		fout.Write(buffer[:n])
	}

	return files
}

// 把多个小的二进制文件合并成一个二进制文件
func mergeBinaryFile(files []string, mergedFile string) {
	os.Remove(mergedFile)
	fout, err := os.OpenFile(mergedFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm) //注意是APPEND模式
	if err != nil {
		log.Panic(err)
	}
	defer fout.Close()

	//遍历每一个小的二进制文件
	for _, file := range files {
		fin, err := os.Open(file)
		if err != nil {
			log.Panic(err)
		}
		defer fin.Close()

		buffer := make([]byte, 1024)
		for { //循环读取一个二进制文件
			n, err := fin.Read(buffer)
			if err != nil {
				if err == io.EOF {
					if n > 0 {
						fout.Write(buffer[:n]) //追加到合并后的大文件里去
					}
				} else {
					log.Println(err)
				}
				break
			} else {
				fout.Write(buffer[:n])
			}
		}
	}
}

func main() {
	utf8 := []byte("老骥伏枥")
	fmt.Println(utf8)
	gbk, _ := Utf8ToGbk(utf8)
	fmt.Println(gbk)
	utf8, _ = GbkToUtf8(gbk) //同一个字符，不同的编码方式将其转为不同的byte数字。
	fmt.Println(utf8)
	fmt.Println(strings.Repeat("*", 50))

	text_file := "io/file_type/a.txt"
	readTextFileByLine(text_file)
	fmt.Println(strings.Repeat("*", 50))
	splitTextFile(text_file)
	fmt.Println(strings.Repeat("*", 50))

	binary_file := "io/file_type/b.png"
	files := splitBinaryFile(binary_file)
	mergedFile := "io/file_type/c.png"
	mergeBinaryFile(files, mergedFile)
}

// go run .\io\file_type\
