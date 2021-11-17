package model

import (
	"encoding/json"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

//认证和授权期间传递的信息
type UserResp struct {
	UserName   string
	Age        int
	UserClaims []Claims //权限列表
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

func GetAdminMsg(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	jsonClaim := claims["claims"].(string) //之所以要解码，是因为在登录的回调函数中存入上下文的时候claims进行了编码
	var userClaims []Claims
	json.Unmarshal([]byte(jsonClaim), &userClaims)
	data := UserResp{UserName: claims["id"].(string), UserClaims: userClaims}
	c.JSON(http.StatusOK, gin.H{
		"code": "success",
		"msg":  "admin访问成功",
		"data": data,
	})
}
