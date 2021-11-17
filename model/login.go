package model

import "github.com/gin-gonic/gin"

func Login(r *gin.Engine, middleware func() gin.HandlerFunc) {
	r.GET("/login", middleware(), func(context *gin.Context) {
		context.JSON(200, gin.H{"msg": "登录中。。"})
	})
}
