package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	db "github.com/thaian1234/simplebank/db/sqlc"
	"github.com/thaian1234/simplebank/utils"
	"os"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymmetricKey:   utils.RandomString(32),
		AccessTokenDuration: time.Minute * 5,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
