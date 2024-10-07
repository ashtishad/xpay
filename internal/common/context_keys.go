package common

type ContextKey string

const (
	AuthorizedUserContextKey ContextKey = "authorizedUser"
	RequestIDContextKey      ContextKey = "requestID"
)
