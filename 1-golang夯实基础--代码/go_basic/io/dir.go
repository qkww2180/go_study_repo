package io

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFile(fileName string) {
	os.Remove(fileName) //先删除，不去理会Remove可能返回的error
	if file, err := os.Create(fileName); err != nil {
		fmt.Printf("create file faied: %v\n", err)
	} else {
		file.Chmod(0o666)                //设置文件权限，八进制
		fmt.Printf("fd=%d\n", file.Fd()) //获取文件描述符file descriptor，这是一个整数
		file.WriteString("多情应笑我\n")
		info, _ := file.Stat()
		fmt.Printf("is dir %t\n", info.IsDir())
		fmt.Printf("modify time %s\n", info.ModTime())
		fmt.Printf("mode %v\n", info.Mode()) //-rw-rw-rw-
		fmt.Printf("file name %s\n", info.Name())
		fmt.Printf("size %dB\n", info.Size()) // 16B
	}

	os.Mkdir("../data/sys", os.ModePerm)          //创建目录并设置权限
	os.MkdirAll("../data/sys/a/b/c", os.ModePerm) //增强版Mkdir，沿途的目录不存在时会一并创建

	// os.Rename("../data/sys/a", "../data/sys/p")       //给文件或目录重命名
	// os.Rename("../data/sys/p/b/c", "../data/sys/p/c") //Rename还可以实现move的功能

	// os.Remove("../data/sys")    //删除文件或目录，目录不为空时才能删除成功
	// os.RemoveAll("../data/sys") //增强版Remove，所有子目录会递归删除
}

// 遍历一个目录
func WalkDir(path string) error {
	//filepath.Walk 会递归地遍历一个目录
	filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.Mode().IsDir() && subPath != path {
			fmt.Printf("path is dir %s\n", subPath)
		} else if info.Mode().IsRegular() {
			fmt.Printf("path is file %s basename %s\n", subPath, info.Name()) //basename表示去除目录，仅保留文件名，subPath里包含了basename
		}
		return nil
	})
	return nil
}
