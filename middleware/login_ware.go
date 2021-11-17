package middleware

import (
	"fmt"
	"gin-login-demo/Utils"
	"github.com/gin-gonic/gin"
)

func LoginMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		//设置context的k-v没有效果的，对于没有用到次中间件的路由，是不会次中间件共享context的
		//context.Set("token","123")
		//context.Set("msg","login success") //没有用到此中间件的路由是无法共享此context的
		fmt.Println("登录中间件执行完了，登录完成")

		Utils.Token = "123" //这样设置才有效
		Utils.Status = "login success"
	}
}
