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
