package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mixdone/uptime-monitoring/internal/models/errs"
)

type tokenService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewTokenService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) TokenService {
	return &tokenService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

func (t *tokenService) Generate(userID int64) (accessToken, refreshToken string, err error) {
	now := time.Now()

	accessClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(t.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString(t.accessSecret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(t.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString(t.refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (t *tokenService) ValidateAccess(tokenStr string) (userID int64, err error) {
	claims, err := t.parseToken(tokenStr, t.accessSecret)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func (t *tokenService) ValidateRefresh(tokenStr string) (userID int64, err error) {
	claims, err := t.parseToken(tokenStr, t.refreshSecret)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func (t *tokenService) parseToken(tokenStr string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.ErrTokenWrongFormat
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errs.ErrTokenExpired
		}
		return nil, errs.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errs.ErrTokenInvalid
	}

	return claims, nil
}
