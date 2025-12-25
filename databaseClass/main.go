package main

import (
	"fmt"

	"volunteer-system/config"
	"volunteer-system/router"
	"volunteer-system/service"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.InitDatabase(); err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	fmt.Println("数据库连接成功")

	// 启动定时任务：检查并关闭过期活动
	go service.CloseExpiredActivities()
	fmt.Println("已启动定时任务：每分钟检查一次过期活动")

	r := gin.Default()

	router.SetupRoutes(r)

	fmt.Println("服务器启动在端口8080")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
