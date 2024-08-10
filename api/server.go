package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/thaian1234/simplebank/api/user"
	db "github.com/thaian1234/simplebank/db/sqlc"
)

const (
	accountsPath  = "/accounts"
	transfersPath = "/transfers"
	usersPath     = "/users"
)

// Server serves HTTP request for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer create a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	gin.ForceConsoleColor()
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	//accounts
	server.registerAccountsRouter(router)
	//transfers
	server.registerTransfersRouter(router)
	//users
	server.registerUsersRouter(router)

	server.registerUserRouter(router)

	server.router = router
	return server
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
	}
}

func (server *Server) registerUserRouter(router *gin.Engine) {
	user := user.NewUserRouter(server.store, "/user")
	user.Routes(router)
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
