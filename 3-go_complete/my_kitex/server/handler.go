package main

import (
	"context"
	math_service "my_kitex/kitex_gen/math_service"
	"time"
)

// MathImpl implements the last service interface defined in the IDL.
type MathImpl struct{}

// Add implements the MathImpl interface.
func (s *MathImpl) Add(ctx context.Context, req *math_service.AddRequest) (resp *math_service.AddResponse, err error) {
	time.Sleep(100 * time.Millisecond)
	resp = &math_service.AddResponse{Sum: req.Left + req.Right}

	// kitex默认为每个handler添加了recover()代码
	// panic("人为制造panic")

	return
}

// Sub implements the MathImpl interface.
func (s *MathImpl) Sub(ctx context.Context, req *math_service.SubRequest) (resp *math_service.SubResponse, err error) {
	time.Sleep(100 * time.Millisecond)
	resp = &math_service.SubResponse{Diff: req.Left - req.Right}
	return
}
