package main

import (
	"log"
	"os"

	"mtmn.top/fish-service/common"
	"mtmn.top/fish-service/routes"
)

func main() {
	common.InitEnv()
	common.InitRedisClient()

	// gin框架启动
	r := routes.SetupRouter()

	// 配置文件读取端口
	port := os.Getenv("PORT")
	log.Println("port config is :" + port)
	ssl := os.Getenv("SSL")
	if ssl == "https" {
		panic(r.RunTLS(":"+port, "./mtmn.top.pem", "./mtmn.top.key"))
	} else {
		panic(r.Run(":" + port))
	}

}
