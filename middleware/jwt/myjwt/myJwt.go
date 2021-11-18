package myjwt

import (
	"encoding/json"
	"fmt"
	"gin-login-demo/Utils"
	"gin-login-demo/middleware/jwt/authorizator"
	"gin-login-demo/model"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"time"
)

var IdentityKey = "id"

func GetGinJWTMiddleware(myAutorizator authorizator.MyAutorizator) (*jwt.GinJWTMiddleware, error) {
	db := Utils.Db
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone", //可以理解成该中间件的名称，用于展示，默认值为gin jwt
		Key:         []byte("secret key"),
		Timeout:     time.Minute * 10, //token过期时间，默认值为time.Hour
		MaxRefresh:  time.Hour,
		IdentityKey: IdentityKey, //身份验证的key值，如果有多个用户的token，会根据这个来识别token,默认值为identity
		PayloadFunc: func(data interface{}) jwt.MapClaims { //登录期间的回调的函数
			//获取权限,从数据库查找对应的权限
			if v, ok := data.(*model.UserResp); ok {
				var loginModel model.LoginModel
				db.Where(&model.LoginModel{Username: v.UserName}).First(&loginModel)
				db.Where(&model.Claims{LoginModelId: loginModel.ID}).First(&v.UserClaims)
				jsonClaim, _ := json.Marshal(v.UserClaims)
				fmt.Println("编码：	", string(jsonClaim))
				return jwt.MapClaims{ //封装身份信息的声明claims,封装到上下文中
					IdentityKey: v.UserName,
					"claims":    string(jsonClaim),
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} { //解析并设置用户身份信息
			//从上下文中提取出之前封装的身份信息的声明claims，是从回调函数返回的
			claims := jwt.ExtractClaims(c)
			jsonClaim := claims["claims"].(string) //claims["claims"]要解码，因为PayloadFunc中进行了编码才存入上下文的
			var userClaims []model.Claims
			json.Unmarshal([]byte(jsonClaim), &userClaims)
			fmt.Println(jsonClaim)
			//返回封装的身份信息,这里返回的信息是后面授权中的信息来源
			return &model.UserResp{
				UserName:   claims[IdentityKey].(string),
				UserClaims: userClaims,
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
		//多态
		Authorizator: myAutorizator.MyAutorizatorMethod, //授权，接收用户信息并编写授权规则，本项目的API权限控制就是通过该函数编写授权规则的

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
		/*
				在调用这个方法ParseToken解析请求中的token时候，会调用jwtFromHeader，
				这个方法会根据空格把token拆成2部分，第一部分就是tokenHeadName,会与这里设置的TokenHeadName比较，不一样就返回错误\
				所以前端需要知道这个TokenHeadName，在携带的时候把这个加在token前才能成功token验证

			所以感觉这样直接用jwt.GinJWTMiddleware还是有局限性的，尝试过重写生成token的方法，太鸡肋了，
			所以还是自己实现生成token比较好！！！！
		*/

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now, //测试或服务器在其他时区可设置该属性，默认值为time.Now
	})
	return authMiddleware, err
}
