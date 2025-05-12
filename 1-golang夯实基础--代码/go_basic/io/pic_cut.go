package io

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func SplitFile(infile string, outDir string, n int) {
	fin, err := os.Open(infile)
	if err != nil {
		log.Panic(err)
	}
	defer fin.Close()

	stat, err := fin.Stat()
	if err != nil {
		log.Panic(err)
	}
	fileSize := stat.Size()
	chunk := fileSize / int64(n)
	if chunk <= 0 {
		panic("file is too small or n is too large")
	}
	for i := 0; i < n; i++ {
		fout, err := os.OpenFile(path.Join(outDir, strconv.Itoa(i)+"_"+path.Base(infile)), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Panic(err)
		}
		need := int(chunk)
		if i == n-1 {
			need = int(fileSize) - (n-1)*int(chunk)
		}
		buffer := make([]byte, need)
		_, err = fin.Read(buffer)
		if err != nil {
			log.Panic(err)
		}

		_, err = fout.Write(buffer)
		if err != nil {
			log.Panic(err)
		}
		fout.Close()
	}
}

func AppendFile(fout *os.File, infile string) {
	fin, err := os.Open(infile)
	if err != nil {
		log.Panic(err)
	}
	defer fin.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := fin.Read(buffer)
		if err != nil {
			if err == io.EOF {
				if n > 0 {
					fout.Write(buffer[:n])
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

// 把dir这个目录下的所有文件合并到mergedFile里去
func MergeFile(dir string, mergedFile string) {
	fout, err := os.OpenFile(mergedFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Panic(err)
	}
	defer fout.Close()

	// ReadDir()得到的子文件/子目录是按文件名排好序的
	if fileInfos, err := os.ReadDir(dir); err != nil {
		log.Panic(err)
	} else {
		for _, fileInfo := range fileInfos {
			if fileInfo.Type().IsRegular() {
				infile := filepath.Join(dir, fileInfo.Name())
				AppendFile(fout, infile)
			}
		}
	}
}
