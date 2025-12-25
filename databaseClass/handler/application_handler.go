package handler

import (
	"net/http"
	"strconv"

	"volunteer-system/model"
	"volunteer-system/service"

	"github.com/gin-gonic/gin"
)

func ApplyActivity(c *gin.Context) {
	idStr := c.Param("id")
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "活动ID格式不正确",
		})
		return
	}

	var req model.ApplyActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}

	application, err := service.ApplyActivity(req.UserID, activityID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "报名成功，等待管理员审核",
		"data":    application,
	})
}

func ListActivityApplications(c *gin.Context) {
	idStr := c.Param("id")
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "活动ID格式不正确",
		})
		return
	}

	apps, err := service.ListActivityApplications(activityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    apps,
	})
}

func ListUserApplications(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "用户ID格式不正确",
		})
		return
	}

	apps, err := service.ListUserApplications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    apps,
	})
}

func UpdateApplicationStatus(c *gin.Context) {
	appIDStr := c.Param("applicationId")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "报名ID格式不正确",
		})
		return
	}

	var req model.UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}

	if err := service.UpdateApplicationStatus(appID, req.Status, req.HandlerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新报名状态成功",
	})
}

// CancelApplication 取消报名
func CancelApplication(c *gin.Context) {
	appIDStr := c.Param("applicationId")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "报名ID格式不正确",
		})
		return
	}

	if err := service.CancelApplication(appID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "取消报名成功",
	})
}
