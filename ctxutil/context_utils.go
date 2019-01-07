package ctxutil

import (
	"context"
	"time"
)

type ctxKey int

const (
	ctxKeyServiceName ctxKey = iota
	ctxKeyCommand
	ctxKeyRequestStartTime
)

func GetServiceName(ctx context.Context) (string, bool) {
	value := ctx.Value(ctxKeyServiceName)
	v, ok := value.(string)
	if !ok {
		return "", false
	}
	return v, true
}

func SetServiceName(ctx context.Context, t string) context.Context {
	return context.WithValue(ctx, ctxKeyServiceName, t)
}

func GetCommand(ctx context.Context) (string, bool) {
	value := ctx.Value(ctxKeyCommand)
	v, ok := value.(string)
	if !ok {
		return "", false
	}
	return v, true
}

func SetCommand(ctx context.Context, t string) context.Context {
	return context.WithValue(ctx, ctxKeyCommand, t)
}

func GetRequestStartTime(ctx context.Context) (time.Time, bool) {
	value := ctx.Value(ctxKeyRequestStartTime)
	t, ok := value.(time.Time)
	if !ok {
		return time.Time{}, false
	}
	return t, true
}

func SetRequestStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, ctxKeyRequestStartTime, t)
}
