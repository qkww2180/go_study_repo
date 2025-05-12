package util

import (
	"os"
	"syscall"
)

// 返回文件的创建时间，最后访问时间，最后修改时间。均为时间戳，精确到秒
func GetFileTime(stat os.FileInfo) (int64, int64, int64) {
	// windows下代码如下
	attr := stat.Sys().(*syscall.Win32FileAttributeData)
	return attr.CreationTime.Nanoseconds() / int64(1e9),
		attr.LastAccessTime.Nanoseconds() / int64(1e9),
		attr.LastWriteTime.Nanoseconds() / int64(1e9)

	// linux环境下代码如下
	//linuxFileAttr := finfo.Sys().(*syscall.Stat_t)
	//return linuxFileAttr.Ctim.Sec,
	// linuxFileAttr.Atim.Sec,
	// linuxFileAttr.Mtim.Sec
}
