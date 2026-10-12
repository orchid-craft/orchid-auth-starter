package constants

type ContextKey string

const (
	RequestContextKey ContextKey = "RequestContextKey"
	UserClaimsContextKey ContextKey = "UserClaimsContextKey"
)

func (c ContextKey) String() string {
	return string(c)
}
