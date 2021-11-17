package model

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// User demo,要返回的用户信息
type UserResp struct {
	UserName string
	Age      int
}

func GetUserInfo(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	fmt.Println(claims)
	userName := claims["id"].(string)
	data := UserResp{UserName: userName}

	c.JSON(http.StatusOK, gin.H{
		"code": "success",
		"msg":  "登录成功",
		"data": data,
	})
}
