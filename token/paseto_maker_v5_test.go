package token

import (
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
	"github.com/thaian1234/simplebank/utils"
	"reflect"
	"testing"
	"time"
)

func TestNewPasetoMakerV5(t *testing.T) {
	type args struct {
		symmetricKey string
	}
	var tests []struct {
		name    string
		args    args
		want    *PasetoMakerV5
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPasetoMakerV5(tt.args.symmetricKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPasetoMakerV5() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPasetoMakerV5() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasetoMakerV5_GenerateToken(t *testing.T) {
	type fields struct {
		paseto       *paseto.V2
		symmetricKey []byte
	}
	type args struct {
		username string
		duration time.Duration
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &PasetoMakerV5{
				paseto:       tt.fields.paseto,
				symmetricKey: tt.fields.symmetricKey,
			}
			got, err := pm.GenerateToken(tt.args.username, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasetoMakerV5_VerifyToken(t *testing.T) {
	type fields struct {
		paseto       *paseto.V2
		symmetricKey []byte
	}
	type args struct {
		tokenString string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    *UserClaims
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &PasetoMakerV5{
				paseto:       tt.fields.paseto,
				symmetricKey: tt.fields.symmetricKey,
			}
			got, err := pm.VerifyToken(tt.args.tokenString)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasetoMakerV5(t *testing.T) {
	tokenString := utils.RandomString(32)
	maker, err := NewPasetoMakerV5(tokenString)
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Hour

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.GenerateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, claims)

	require.NotZero(t, claims.ID)
	require.Equal(t, username, claims.Username)
	require.WithinDuration(t, issuedAt, claims.IssuedAt.Time, time.Second)
	require.WithinDuration(t, expiredAt, claims.ExpiresAt.Time, time.Second)
}
