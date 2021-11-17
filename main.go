package main

import (
	"gin-login-demo/routers"
	"net/http"
)

func main() {
	r := routers.InitRoutes()
	server := &http.Server{Addr: "127.0.0.1:8088", Handler: r}
	server.ListenAndServe()
	//r.Run("127.0.0.1:8088")
}
