package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/thaian1234/simplebank/db/sqlc"
	"github.com/thaian1234/simplebank/token"
	"github.com/thaian1234/simplebank/utils"
)

const (
	accountsPath  = "/accounts"
	transfersPath = "/transfers"
	usersPath     = "/users"
)

// Server serves HTTP request for our banking service
type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.MakerV5
	config     utils.Config
}

// NewServer create a new HTTP server and setup routing
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMakerV5(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
	gin.ForceConsoleColor()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
		}
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	//accounts
	server.registerAccountsRouter(router)
	//transfers
	server.registerTransfersRouter(router)
	//users
	server.registerUsersRouter(router)

	server.router = router
}

func (server *Server) registerAccountsRouter(router *gin.Engine) {
	accounts := router.Group(accountsPath)
	{
		accounts.POST("/", server.createAccount)
		accounts.GET("/:id", server.getAccount)
		accounts.GET("/", server.getListAccounts)
	}
}

func (server *Server) registerTransfersRouter(router *gin.Engine) {
	transfers := router.Group(transfersPath)
	{
		transfers.POST("/", server.createTransfer)
	}
}

func (server *Server) registerUsersRouter(router *gin.Engine) {
	users := router.Group(usersPath)
	{
		users.POST("/", server.createUser)
		users.POST("/login", server.loginUser)
	}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func msgResponse(msg string) gin.H {
	return gin.H{
		"msg": msg,
	}
}
