package model

import (
	"fmt"
	"gin-login-demo/Utils"
	"github.com/gin-gonic/gin"
)

func OrderGet(r *gin.Engine) {
	r.GET("/order", func(context *gin.Context) {
		//token,_ := context.Get("token") //这样是取不到token的，因为loginware中间件没作用在此路由上，context不是同一个
		token := Utils.Token
		fmt.Println(token)
		context.JSON(200, gin.H{"token": token, "msg": "orders"})
	})
}
