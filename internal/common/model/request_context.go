package model

import (
	"context"

	"github.com/yudhiana/web-api/constants"
)

// RequestContext contains common request metadata extracted from headers
type RequestContext struct {
	RequestID string
}

// WithRequestContext adds RequestContext to the given context
func WithRequestContext(ctx context.Context, key any, value any) context.Context {
	return context.WithValue(ctx, value, key)
}

func SetRequestContext(ctx context.Context, reqCtx *RequestContext) context.Context {
	return context.WithValue(ctx, constants.RequestContextKey, reqCtx)
}

// GetRequestContext retrieves RequestContext from the given context
func GetRequestContext(ctx context.Context) (*RequestContext, bool) {
	reqCtx, ok := ctx.Value(constants.RequestContextKey).(*RequestContext)
	return reqCtx, ok
}
