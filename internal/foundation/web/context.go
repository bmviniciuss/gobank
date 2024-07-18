package web

import "context"

type ctxKey int

const (
	reqIDKey ctxKey = iota + 1
)

func setRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, reqIDKey, id)
}

func GetRequestID(ctx context.Context) string {
	v, ok := ctx.Value(reqIDKey).(string)
	if !ok {
		return "00000000-0000-0000-0000-000000000000"
	}
	return v
}
