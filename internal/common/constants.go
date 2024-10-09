package common

const (
	AppEnvDev        = "dev"
	AppEnvProduction = "production"
	AppEnvStaging    = "staging"

	RequestIDHeader = "X-Request-ID"
	RequestIDKey    = "requestID"

	DBTSLayout       = "time.RFC3339"
	CardExpiryLayout = "01/06" // MM/YY

	DBColumnID       = "id"
	DBColumnUUID     = "uuid"
	DBColumnUserID   = "user_id"
	DBColumnEmail    = "email"
	DBColumnWalletID = "wallet_id"
)
