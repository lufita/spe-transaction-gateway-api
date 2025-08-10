package lookup

import "os"

const (
	SUCCESS_CODE          = "00"
	UNAUTHORIZED_CODE     = "401"
	INTERNAL_SERVER_ERROR = "500"
	BAD_REQUEST           = "Bad Request"
)

var (
	RedisHost              = os.Getenv("REDIS_HOST")
	RedisPort              = os.Getenv("REDIS_PORT")
	DbConnString           = os.Getenv("DB_CONNECTION_STRING")
	NotificationSignSecret = os.Getenv("NOTIF_SIGNING_SECRET")
	InquirySignSecret      = os.Getenv("INQUIRY_SIGN_SECRET")
	JWTTokenSecret         = os.Getenv("AUTHORIZATION_TOKEN_SECRET")
	RegisterSecret         = os.Getenv("REGISTER_SECRET")
)
