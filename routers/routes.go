package routers

import (
	"gin-login-demo/Utils"
	"gin-login-demo/model"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var IdentityKey = "id"

func InitRoutes() *gin.Engine {
	db := Utils.Db
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//model.Login(r, middleware.LoginMiddleware)

	// 获取 jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone", //可以理解成该中间件的名称，用于展示，默认值为gin jwt
		Key:         []byte("secret key"),
		Timeout:     time.Minute * 10, //token过期时间，默认值为time.Hour
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey, //身份验证的key值，如果有多个用户的token，会根据这个来识别token,默认值为identity
		PayloadFunc: func(data interface{}) jwt.MapClaims { //登录期间的回调的函数
			if v, ok := data.(*model.UserResp); ok {
				return jwt.MapClaims{
					IdentityKey: v.UserName, //封装身份信息的声明claims
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} { //解析并设置用户身份信息
			//从上下文中提取出之前封装的身份信息的声明claims
			claims := jwt.ExtractClaims(c)
			//返回封装的身份信息
			return &model.UserResp{
				UserName: claims[IdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) { //校验登录信息的回调信息,认证
			var loginVals model.LoginModel //登录信息
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			//从数据库获取并校验
			lo := model.LoginModel{}
			db.Where(&model.LoginModel{Username: username, Password: password}).First(&lo)
			if lo.ID > 0 {
				return &model.UserResp{
					UserName: lo.Username,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool { //授权，接收用户信息并编写授权规则，本项目的API权限控制就是通过该函数编写授权规则的
			//授权逻辑
			//这里假装一下，给名字是wb的授权
			if v, ok := data.(*model.UserResp); ok && v.UserName == "wb" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) { //处理不进行授权的逻辑
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt", //token检索模式，用于提取token，默认值为header:Authorization
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer", //token在请求头时的名称，默认值为Bearer

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now, //测试或服务器在其他时区可设置该属性，默认值为time.Now
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// When you use jwt.New(), the function is already automatically called for checking,
	// which means you don't need to call it again.
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

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

	return r
}
