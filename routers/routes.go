package routers

import (
	"gin-login-demo/middleware"
	"gin-login-demo/model"
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	model.Login(r, middleware.LoginMiddleware)
	model.UserGet(r)
	model.OrderGet(r)
	return r
}
