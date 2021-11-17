package authorizator

import (
	"gin-login-demo/model"
	"github.com/gin-gonic/gin"
)

//自定义授权器，以满足不同的路由授权

type MyAutorizator interface { //自定义授权器接口
	MyAutorizatorMethod(data interface{}, c *gin.Context) bool
} //jwt中的授权器是func(data interface{}, c *gin.Context) bool 类型的

//之后就是多态实现此接口，再把想要的实例传入jwt中间件创建时的授权器参数中
type AdminAutorizator struct{} //定义一个admin的授权器

func (*AdminAutorizator) MyAutorizatorMethod(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*model.UserResp); ok { //类型断言
		for _, claims := range v.UserClaims {
			if claims.Type == "role" && claims.Value == "admin" {
				return true
			}
		}
	}
	return false
}

type AllAuthorizator struct{} //定义一个全权限的授权器

func (*AllAuthorizator) MyAutorizatorMethod(data interface{}, c *gin.Context) bool {
	return true
}
