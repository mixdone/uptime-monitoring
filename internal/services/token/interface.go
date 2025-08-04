package token

type TokenService interface {
	Generate(userID int64) (accessToken, refreshToken string, err error)
	ValidateAccess(tokenStr string) (userID int64, err error)
	ValidateRefresh(tokenStr string) (userID int64, err error)
}
