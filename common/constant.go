package common

import (
	"errors"
)

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusError   = "error"

	Timezone        = "TIMEZONE"
	DefaultTimeZone = "Asia/Jakarta"

	ConsulWatchInterval       = "CONSUL_WATCH_INTERVAL_SECONDS"
	DefaultLoadConsulInterval = 30

	XRequestIdHeader = "x-request-id"

	StatusUserActive   = "active"
	StatusUserInactive = "inactive"
)

var (
	ErrSQLQueryBuilder      = errors.New("error query builder")
	ErrSQLExec              = errors.New("error sql exec")
	ErrDataNotFound         = errors.New("data not found")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrEmailAlreadyTaken    = errors.New("email already taken")
	ErrFailHashPassword     = errors.New("failed to hash password")
	ErrUserAlreadyDeleted   = errors.New("user already deleted")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrInvalidToken         = errors.New("invalid token")
)
