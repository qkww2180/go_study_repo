package j_io

import (
	"bufio"
	"fmt"
	"os"
)

func WriteFile() {
	// 如果使用go test，则相对路径是相对于xxx_test.go文件的路径。
	// 如果使用go run或编译后直接运行，则相对路径是相对于执行命令时所在的路径。
	if fout, err := os.OpenFile("../z_data/verse.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666); err != nil {
		fmt.Printf("open file failed: %s\n", err.Error())
	} else {
		defer fout.Close()
		fout.WriteString("纳兰性德\n")
		fout.WriteString("明月多情应笑我")
		fout.WriteString("\n")
		fout.Write([]byte("笑我如今"))
		fout.WriteString("\n")
	}
}

func WriteFileWithBuffer() {
	if fout, err := os.OpenFile("../z_data/verse.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666); err != nil {
		fmt.Printf("open file failed: %s\n", err.Error())
	} else {
		defer fout.Close()
		writer := bufio.NewWriter(fout)
		writer.WriteString("纳兰性德\n")
		writer.Write([]byte("明月多情应笑我\n"))
		writer.WriteString("笑我如今\n")
		writer.Flush()
	}
}
