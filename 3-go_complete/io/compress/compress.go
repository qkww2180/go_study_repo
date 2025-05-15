package compress

import (
	"compress/gzip" //compress下还有zlib,bzip2,lzw等其他压缩算法
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

type Writer interface {
	Write([]byte) (int, error)
	Close() error
}

const (
	COMPRESS_ALGO = iota //0
	ZLIB                 //1
	GZIP                 //2
)

func getCompressWriter(compressAlgo int, file *os.File) (Writer, error) { //返回一个接口
	switch compressAlgo {
	case ZLIB:
		return zlib.NewWriter(file), nil
	case GZIP:
		return gzip.NewWriter(file), nil
	default:
		return nil, fmt.Errorf("INVALID compress algo %d", compressAlgo)
	}
}

func getDecompressReader(compressAlgo int, file *os.File) (io.ReadCloser, error) {
	switch compressAlgo {
	case ZLIB:
		return zlib.NewReader(file)
	case GZIP:
		return gzip.NewReader(file)
	default:
		return nil, fmt.Errorf("INVALID decompress algo %d", compressAlgo)
	}
}

// 压缩
func compress(compressAlgo int, inFile, outFile string) {
	fin, err := os.Open(inFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	stat, _ := fin.Stat()
	fmt.Printf("压缩前文件大小 %dB\n", stat.Size())

	fout, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	writer, err := getCompressWriter(compressAlgo, fout) //压缩写入
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// bs := make([]byte, 1024)
	// for {
	// 	n, err := fin.Read(bs)
	// 	if err != nil {
	// 		if err == j_io.EOF {
	// 			break
	// 		} else {
	// 			fmt.Println(err)
	// 		}
	// 	} else {
	// 		writer.Write(bs[:n])
	// 	}
	// }
	io.Copy(writer, fin) //把一个流拷贝到另外一个流

	writer.Close()
	fout.Close()
	fin.Close()
}

// 解压
func decompress(compressAlgo int, inFile, outFile string) {
	fin, err := os.Open(inFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	stat, _ := fin.Stat()
	fmt.Printf("压缩后文件大小 %dB\n", stat.Size())

	reader, err := getDecompressReader(compressAlgo, fin) //解压
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fout, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(fout, reader) //把一个流拷贝到另外一个流
	reader.Close()
	fin.Close()
	fout.Close()
}
