package handler

import (
	"net/http"

	"volunteer-system/service"

	"github.com/gin-gonic/gin"
)

// GetStatistics 获取系统统计数据
func GetStatistics(c *gin.Context) {
	stats, err := service.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// 获取部门统计数据 (聚合函数+GROUP BY)
func GetDeptStatistics(c *gin.Context) {
	stats, err := service.GetDeptStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// 获取分类统计数据 (聚合函数+GROUP BY)
func GetCategoryStatistics(c *gin.Context) {
	stats, err := service.GetCategoryStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// 获取用户活跃度统计 (聚合函数+GROUP BY+HAVING)
func GetUserActivityStatistics(c *gin.Context) {
	stats, err := service.GetUserActivityStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// 获取活动热度排行 (聚合函数+GROUP BY)
func GetActivityPopularity(c *gin.Context) {
	stats, err := service.GetActivityPopularity()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// 获取管理员创建活动统计 (聚合函数+GROUP BY)
func GetAdminCreationStatistics(c *gin.Context) {
	stats, err := service.GetAdminCreationStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// 获取即将开始的活动 (日期时间函数)
func GetUpcomingActivities(c *gin.Context) {
	stats, err := service.GetUpcomingActivities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// 获取热门活动 (关联子查询)
func GetPopularActivities(c *gin.Context) {
	stats, err := service.GetPopularActivities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetOmnipotentVolunteers 获取全能志愿者列表（除法查询）
func GetOmnipotentVolunteers(c *gin.Context) {
	volunteers, err := service.GetOmnipotentVolunteers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    volunteers,
	})
}

// GetUserApplicationInfo 获取所有用户及其报名情况（外连接）
func GetUserApplicationInfo(c *gin.Context) {
	info, err := service.GetUserApplicationInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}

// GetDeptActivityInfo 获取所有部门及其活动数（外连接）
func GetDeptActivityInfo(c *gin.Context) {
	info, err := service.GetDeptActivityInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}
