package raft

import (
	rand "math/rand/v2"
	"net/url"
	"path"
	"time"
)

// 返回1-2倍的timeout
func randomTimeout(timeout time.Duration) <-chan time.Time {
	if timeout == 0 {
		return nil
	}
	delta := time.Duration(rand.Int64()) % timeout
	return time.After(timeout + delta)
}

// 把thePath拼接到connString后面
func joinUrlPath(connString, thePath string) string {
	u, err := url.Parse(connString)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, thePath)
	return u.String()
}
