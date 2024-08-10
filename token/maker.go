package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)

	//VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}

type MakerV5 interface {
	GenerateToken(username string, duration time.Duration) (string, error)

	//VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*JWTData, error)
}
