package routers

import (
	"gin-login-demo/middleware/jwt/authorizator"
	"gin-login-demo/middleware/jwt/myjwt"
	"gin-login-demo/model"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var IdentityKey = "id"

func InitRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 获取 jwt middleware
	authMiddleware, err := myjwt.GetGinJWTMiddleware(&authorizator.AllAuthorizator{}) //传入自定义的授权器
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	//errInit := authMiddleware.MiddlewareInit()

	//if errInit != nil {
	//	log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	//}

	//使用jwt中间件
	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	api := r.Group("/user")
	api.Use(authMiddleware.MiddlewareFunc()) //对所有user组的路由使用authMiddleware中间件
	{
		api.GET("/info", model.GetUserInfo)
	}

	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc()) //对所有auth组的路由使用authMiddleware中间件
	{
		// Refresh time can be longer than token timeout
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	admin_authMiddleware, err := myjwt.GetGinJWTMiddleware(&authorizator.AdminAutorizator{}) //传入自定义的授权器
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	admin := r.Group("/admin")
	admin.Use(admin_authMiddleware.MiddlewareFunc())
	{
		admin.GET("/adminMsg", model.GetAdminMsg)
	}

	return r
}
