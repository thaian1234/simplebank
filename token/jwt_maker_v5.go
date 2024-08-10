package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMakerV5 struct {
	secretKey []byte
}

func NewJWTMakerV5(secretKey string) (MakerV5, error) {
	if len(secretKey) < minSercretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSercretKeySize)
	}
	return &JWTMakerV5{
		secretKey: []byte(secretKey),
	}, nil
}

func (maker *JWTMakerV5) GenerateToken(username string, duration time.Duration) (string, error) {
	jwtData, err := NewJWTData(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtData)
	signedToken, err := jwtToken.SignedString(maker.secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (maker *JWTMakerV5) VerifyToken(tokenString string) (*JWTData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTData{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return maker.secretKey, nil
	})

	if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, jwt.ErrSignatureInvalid
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		// Invalid signature
		return nil, jwt.ErrTokenSignatureInvalid
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		return nil, jwt.ErrTokenExpired
	}

	claims, ok := token.Claims.(*JWTData)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
