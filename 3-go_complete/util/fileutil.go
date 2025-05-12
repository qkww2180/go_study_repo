package util

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	BlankReg = regexp.MustCompile(`\s+`)
)

// ReadAllLines 读取文件中的所有行
func ReadAllLines(infile string) []string {
	lines := make([]string, 0)
	f, err := os.Open(infile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		lines = append(lines, strings.TrimSpace(line))
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
	}
	return lines
}

func Md5SumFile(file string) (value [md5.Size]byte, err error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return
	}
	value = md5.Sum(data)
	return
}

func byteAdd(arr, brr [md5.Size]byte) [md5.Size]byte {
	crr := [md5.Size]byte{}
	for i := 0; i < md5.Size; i++ {
		crr[i] = arr[i] + brr[i] //如果溢出就不管了
	}
	return crr
}

func IterFolder(folder string) []string {
	FileList := []string{}
	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.Mode().IsDir() && path != folder {
			FileList = append(FileList, IterFolder(path)...) //闭包中使用外部变量时传的是引用。golang的传参数全部是传值，没有传引用的
		} else if info.Mode().IsRegular() {
			FileList = append(FileList, path)
		}
		return nil
	})

	dup := make(map[string]bool)
	for _, ele := range FileList {
		dup[ele] = true
	}
	FileList = []string{}
	for k := range dup {
		FileList = append(FileList, k)
	}
	return FileList
}

// 判断文件或文件夹是否存在。
// 如果返回的错误为nil,说明文件或文件夹存在。
// 如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在。
// 如果返回的错误为其它类型,则不确定是否在存在。
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 解压
func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		err = os.MkdirAll(path.Dir(filename), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}
