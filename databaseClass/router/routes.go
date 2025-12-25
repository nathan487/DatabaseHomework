package router

import (
	"net/http"
	"time"

	"volunteer-system/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// 提供静态文件服务
	r.Static("/static", "./")

	// 提供前端页面
	r.GET("/", func(c *gin.Context) {
		// 检查是否已登录，决定显示首页还是登录页
		c.File("./index.html")
	})
	r.GET("/login", func(c *gin.Context) {
		c.File("./login.html")
	})
	r.GET("/test", func(c *gin.Context) {
		c.File("./test.html")
	})

	// User routes
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)

	// Activity routes
	activityGroup := r.Group("/activities")
	{
		activityGroup.GET("", handler.ListActivities)
		activityGroup.GET("/search", handler.SearchActivities)
		activityGroup.GET("/upcoming", handler.GetUpcomingActivities)
		activityGroup.GET("/popular", handler.GetPopularActivities)
		activityGroup.GET("/available", handler.GetAvailableActivities)
		activityGroup.GET("/:id", handler.GetActivityDetail)
		activityGroup.POST("", handler.CreateActivity)
		activityGroup.PUT("/:id", handler.UpdateActivity)
		activityGroup.DELETE("/:id", handler.DeleteActivity)
		activityGroup.POST("/:id/apply", handler.ApplyActivity)
		activityGroup.GET("/:id/applications", handler.ListActivityApplications)
	}

	// Application routes
	r.GET("/users/:userId/applications", handler.ListUserApplications)
	r.POST("/applications/:applicationId/status", handler.UpdateApplicationStatus)
	r.DELETE("/applications/:applicationId", handler.CancelApplication)

	// Statistics routes
	r.GET("/statistics", handler.GetStatistics)
	r.GET("/statistics/departments", handler.GetDeptStatistics)
	r.GET("/statistics/categories", handler.GetCategoryStatistics)
	r.GET("/statistics/users", handler.GetUserActivityStatistics)
	r.GET("/statistics/activities/popularity", handler.GetActivityPopularity)
	r.GET("/statistics/admins", handler.GetAdminCreationStatistics)
	r.GET("/statistics/omnipotent-volunteers", handler.GetOmnipotentVolunteers)
}
