package common

import (
	"context"
	"net/http"

	"github.com/yudhiana/web-api/constants"
	modelCommon "github.com/yudhiana/web-api/internal/common/model"
)

// ExtractRequestContext extracts common headers from http.Request and creates RequestContext
func ExtractRequestContext(r *http.Request) *modelCommon.RequestContext {
	return &modelCommon.RequestContext{
		RequestID: r.Header.Get(constants.HeaderRequestID),
	}
}

func SetRequestContext(ctx context.Context, r *http.Request) context.Context {
	reqCtx := ExtractRequestContext(r)
	return modelCommon.SetRequestContext(ctx, reqCtx)
}

// GetRequestIDFromContext extracts request ID from context
func GetRequestIDFromContext(ctx context.Context) string {
	if reqCtx, ok := modelCommon.GetRequestContext(ctx); ok {
		return reqCtx.RequestID
	}
	return ""
}
