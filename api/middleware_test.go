package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/thaian1234/simplebank/token"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func addAuthorization(t *testing.T, req *http.Request, tokenMaker token.MakerV5, authorizationType string, username string, duration time.Duration) {
	generatedToken, err := tokenMaker.GenerateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, generatedToken)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, generatedToken)
	req.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, req *http.Request, tokenMaker token.MakerV5)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{{
		name: "OK",
		setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.MakerV5) {
			addAuthorization(t, req, tokenMaker, authorizationTypeBearer, "user", time.Minute*1)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
		},
	}, {
		name: "NoAuthorization",
		setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.MakerV5) {
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusUnauthorized, recorder.Code)
		},
	}, {
		name: "UnsupportedAuthorization",
		setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.MakerV5) {
			addAuthorization(t, req, tokenMaker, "unsupported", "user", time.Minute)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusUnauthorized, recorder.Code)
		},
	}, {
		name: "UnsupportedAuthorization",
		setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.MakerV5) {
			addAuthorization(t, req, tokenMaker, "", "user", time.Minute)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusUnauthorized, recorder.Code)
		},
	}, {
		name: "ExpiredToken",
		setupAuth: func(t *testing.T, req *http.Request, tokenMaker token.MakerV5) {
			addAuthorization(t, req, tokenMaker, authorizationTypeBearer, "user", -time.Minute)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusUnauthorized, recorder.Code)
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			server := newTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(authPath, authMiddleware(server.tokenMaker), func(context *gin.Context) {
				context.JSON(http.StatusOK, gin.H{})
			})

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)
			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})

	}
}
