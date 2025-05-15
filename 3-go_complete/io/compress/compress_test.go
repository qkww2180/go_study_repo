package compress

import (
	"dqq/io/util"
	"testing"
)

func TestGzip(t *testing.T) {
	flatFile1 := util.RootPath + "go.mod"
	flatFile2 := util.RootPath + "go.mod.copy"
	compressFile := util.RootPath + "go.mod.gzip"
	compress(GZIP, flatFile1, compressFile)
	decompress(GZIP, compressFile, flatFile2)
}

func TestZlib(t *testing.T) {
	flatFile1 := util.RootPath + "go.mod"
	flatFile2 := util.RootPath + "go.mod.copy"
	compressFile := util.RootPath + "go.mod.zlib"
	compress(ZLIB, flatFile1, compressFile)
	decompress(ZLIB, compressFile, flatFile2)
}

//go test -v .\j_io\compress\ -run=^TestGzip$ -count=1
//go test -v .\j_io\compress\ -run=^TestZlib$ -count=1
