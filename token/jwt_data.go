package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTData struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

func NewJWTData(username string, duration time.Duration) (*JWTData, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	jwtData := &JWTData{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			ID:        tokenID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
		Username: username,
	}
	return jwtData, nil
}

func (jd *JWTData) Validate() error {
	if time.Now().After(jd.ExpiresAt.Time) {
		return ErrExpiredToken
	}
	return nil
}
