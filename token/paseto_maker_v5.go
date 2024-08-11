package token

import (
	"fmt"
	"github.com/aead/chacha20poly1305"
	"github.com/golang-jwt/jwt/v5"
	"github.com/o1egl/paseto"
	"time"
)

type PasetoMakerV5 struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMakerV5(symmetricKey string) (*PasetoMakerV5, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMakerV5{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte((symmetricKey)),
	}
	return maker, nil
}

func (pm *PasetoMakerV5) GenerateToken(username string, duration time.Duration) (string, error) {
	claims, err := NewUserClaims(username, duration)
	if err != nil {
		return "", err
	}
	return pm.paseto.Encrypt(pm.symmetricKey, claims, nil)
}

func (pm *PasetoMakerV5) VerifyToken(tokenString string) (*UserClaims, error) {
	claims := &UserClaims{}
	if err := pm.paseto.Decrypt(tokenString, pm.symmetricKey, claims, nil); err != nil {
		return nil, jwt.ErrTokenInvalidId
	}
	if err := claims.Validate(); err != nil {
		return nil, err
	}
	return claims, nil
}
