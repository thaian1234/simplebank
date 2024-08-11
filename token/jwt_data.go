package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

func NewUserClaims(username string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	jwtData := &UserClaims{
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

func (jd *UserClaims) Validate() error {
	if time.Now().After(jd.ExpiresAt.Time) {
		return jwt.ErrTokenExpired
	}
	return nil
}
