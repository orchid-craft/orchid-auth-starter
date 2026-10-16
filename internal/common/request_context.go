package common

import (
	"context"
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/olympus-forge/logos"
	"github.com/orchid-craft/orchid-auth-starter/constants"
	modelCommon "github.com/orchid-craft/orchid-auth-starter/internal/common/model"
)

// ExtractRequestContext extracts common headers from http.Request and creates RequestContext
func ExtractRequestContext(r *http.Request) *modelCommon.RequestContext {
	reqID := r.Header.Get(constants.HeaderRequestID)
	if reqID == "" {
		reqID = chiMiddleware.GetReqID(r.Context())
	}
	return &modelCommon.RequestContext{
		RequestID: reqID,
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

// LogInfoWithContext menulis log INFO dengan menyertakan request_id dari context jika ada (P2-4)
func LogInfoWithContext(ctx context.Context, logger *logos.LogEntry, msg string, args ...any) {
	if reqID := GetRequestIDFromContext(ctx); reqID != "" {
		args = append(args, "request_id", reqID)
	}
	logger.Info(msg, args...)
}

// LogErrorWithContext menulis log ERROR dengan menyertakan request_id dari context jika ada (P2-4)
func LogErrorWithContext(ctx context.Context, logger *logos.LogEntry, msg string, args ...any) {
	if reqID := GetRequestIDFromContext(ctx); reqID != "" {
		args = append(args, "request_id", reqID)
	}
	logger.Error(msg, args...)
}
