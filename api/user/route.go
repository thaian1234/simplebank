package user

import (
	"github.com/gin-gonic/gin"
	db "github.com/thaian1234/simplebank/db/sqlc"
)

type UserRouter struct {
	store db.Store
	path  string
}

func NewUserRouter(store db.Store, path string) *UserRouter {
	return &UserRouter{
		store: store,
		path:  path,
	}
}

func (u *UserRouter) Routes(router *gin.Engine) {
	user := router.Group(u.path)
	{
		user.POST("/", u.createUserHandler)
	}
}
