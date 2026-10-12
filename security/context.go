package security

import (
	"context"

	"github.com/yudhiana/web-api/constants"
)

// SetClaimsContext menyimpan claims JWT pengguna ke dalam context
func SetClaimsContext(ctx context.Context, claims *AccessTokenClaims) context.Context {
	return context.WithValue(ctx, constants.UserClaimsContextKey, claims)
}

// GetClaimsContext mengambil claims JWT pengguna dari context
func GetClaimsContext(ctx context.Context) (*AccessTokenClaims, bool) {
	claims, ok := ctx.Value(constants.UserClaimsContextKey).(*AccessTokenClaims)
	return claims, ok
}

// GetUserIDFromContext helper untuk mengambil UserID secara langsung
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	claims, ok := GetClaimsContext(ctx)
	if !ok || claims == nil {
		return "", false
	}
	return claims.UserID, true
}
