package main

import (
	"fmt"
	"gin-login-demo/Utils"
	"gin-login-demo/routers"
	"github.com/jinzhu/gorm"
	"net/http"
)

var db *gorm.DB

func init() {
	//用户名:密码@tcp(数据库ip或域名:端口)/数据库名称?charset=数据库编码&parseTime=True&loc=Local
	var err error
	db, err = gorm.Open("mysql", "root:12345@tcp(127.0.0.1:3306)/db_01?"+
		"charset=utf8&parseTime=True&loc=Local")
	//有点像go的数据库包一样，使用open方法来两将诶数据库
	if err != nil {
		panic(err)
	}
	// fmt.Printf("%#v\n", db)

	db.LogMode(true) //开启日志打印
	//db.SetLogger(log.New(os.Stdout, "\r\n", 0)) //设置日志格式

	//配置连接池
	db.DB().SetMaxIdleConns(10)  //最大空闲连接池数
	db.DB().SetMaxOpenConns(100) //数据库打开的最大连接数

	Utils.Db = db

	fmt.Println(Utils.Db)
}

func main() {
	r := routers.InitRoutes()
	server := &http.Server{Addr: "127.0.0.1:8088", Handler: r}
	server.ListenAndServe()
	//r.Run("127.0.0.1:8088")
}
