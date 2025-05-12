package io

import (
	"context"
	"log/slog"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

type SlogContextHandler struct {
	slog.Handler
}

// 添加上下文属性到 Record 中，然后调用底层的 handler
func (h SlogContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			r.AddAttrs(v)
		}
	}

	return h.Handler.Handle(ctx, r)
}

// AppendCtx 将 slog 属性添加到提供的上下文中，
// 以便在使用此类上下文创建的任何 Record 中都会包含该属性
func AppendCtx(parent context.Context, attr slog.Attr) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(parent, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(parent, slogFields, v)
}
