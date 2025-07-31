package constants

import "time"

const (
	RefreshTokenTTL = 7 * 24 * time.Hour
	AccessTokenTTL  = 15 * time.Minute
)
