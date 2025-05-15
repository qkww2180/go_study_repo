package j_io_test

import (
	"dqq/go/basic/j_io"
	"testing"
)

func TestWriteFile(t *testing.T) {
	j_io.WriteFile()
}

func TestReadFile(t *testing.T) {
	j_io.ReadFile()
}

func TestWriteFileWithBuffer(t *testing.T) {
	j_io.WriteFileWithBuffer()
}

func TestReadFileWithBuffer(t *testing.T) {
	j_io.ReadFileWithBuffer()
}

func TestSplitFile(t *testing.T) {
	imgFile := "../z_img/大乔乔好课.png"
	j_io.SplitFile(imgFile, "../z_img/图像分割", 4)
}

func TestMergeFile(t *testing.T) {
	j_io.MergeFile("../z_img/图像分割", "../z_img/图像合并.png")
}

func TestCreateFile(t *testing.T) {
	j_io.CreateFile("../z_data/poem.txt")
}

func TestWalkDir(t *testing.T) {
	j_io.WalkDir("../z_data")
}

func TestCompress(t *testing.T) {
	j_io.Compress("../z_img/大乔乔好课.png", "../z_img/大乔乔好课.png.gzip", j_io.GZIP)
	j_io.Decompress("../z_img/大乔乔好课.png.gzip", "../z_data/大乔乔好课.png", j_io.GZIP)
}

func TestLog(t *testing.T) {
	logger := j_io.NewLogger("../z_data/biz.log")
	j_io.Log(logger)
}

func TestSLog(t *testing.T) {
	logger := j_io.NewSLogger("../z_data/sbiz.log")
	j_io.SLog(logger)
}

func TestSysCall(t *testing.T) {
	j_io.SysCall()
}

func TestJson(t *testing.T) {
	j_io.JsonSerialize()
}

func TestRegex(t *testing.T) {
	j_io.UseRegex()
}

// go test -v ./j_io -run=^TestWriteFile$ -count=1
// go test -v ./j_io -run=^TestReadFile$ -count=1
// go test -v ./j_io -run=^TestWriteFileWithBuffer$ -count=1
// go test -v ./j_io -run=^TestReadFileWithBuffer$ -count=1
// go test -v ./j_io -run=^TestCreateFile$ -count=1
// go test -v ./j_io -run=^TestWalkDir$ -count=1
// go test -v ./j_io -run=^TestSplitFile$ -count=1
// go test -v ./j_io -run=^TestMergeFile$ -count=1
// go test -v ./j_io -run=^TestJson$ -count=1
// go test -v ./j_io -run=^TestCompress$ -count=1
// go test -v ./j_io -run=^TestLog$ -count=1
// go test -v ./j_io -run=^TestSLog$ -count=1
// go test -v ./j_io -run=^TestSysCall$ -count=1
// go test -v ./j_io -run=^TestRegex$ -count=1
