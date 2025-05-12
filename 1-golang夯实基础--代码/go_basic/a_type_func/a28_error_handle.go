package main

import (
	"errors"
	"fmt"
)

var (
	ErrNotExists      = errors.New("file or directory not exists")
	ErrNotPermitted   = errors.New("not permitted")
	ErrEmptyDirectory = errors.New("could not delete empty directory")
)

// 调用链  A->B->C  A处于最上游

func deleteDir(dir string) error {
	var (
		dirNotExists bool
		forbidden    bool
		emptyDir     bool
	)
	if dirNotExists {
		fmt.Printf("%s 不存在\n", dir)
		return nil
		// return ErrNotExists //可能不需要把error抛给上游，内部把error消化掉即可
	}
	if forbidden {
		return ErrNotPermitted
	}
	if emptyDir {
		return ErrEmptyDirectory
	}
	return nil
}

func scheduleTask() {
	dir := "/a/b/c"
	err := deleteDir(dir)
	if err != nil {
		if err == ErrNotPermitted {
			// if errors.Is(err, ErrNotPermitted) {  //errors.Is和==并不等价
			err = deleteDir(dir) //重试，或者通过其他途径完成工作
		}
		if err != nil {
			fmt.Printf("删除dir %s 失败 %s\n", dir, err)
		}
	}
}

func main28() {
	scheduleTask()
}
