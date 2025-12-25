package handler

import (
	"net/http"
	"strconv"

	"volunteer-system/model"
	"volunteer-system/service"

	"github.com/gin-gonic/gin"
)

func ListActivities(c *gin.Context) {
	var deptID, categoryID *int

	if deptIDStr := c.Query("dept_id"); deptIDStr != "" {
		if id, err := strconv.Atoi(deptIDStr); err == nil {
			deptID = &id
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if id, err := strconv.Atoi(categoryIDStr); err == nil {
			categoryID = &id
		}
	}

	activities, err := service.ListActivities(deptID, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    activities,
	})
}

func CreateActivity(c *gin.Context) {
	var req model.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}

	activity, err := service.CreateActivity(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "活动创建成功",
		"data":    activity,
	})
}

func UpdateActivity(c *gin.Context) {
	idStr := c.Param("id")
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "活动ID格式不正确",
		})
		return
	}

	var req model.UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}

	activity, err := service.UpdateActivity(activityID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "活动更新成功",
		"data":    activity,
	})
}

func DeleteActivity(c *gin.Context) {
	idStr := c.Param("id")
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "活动ID格式不正确",
		})
		return
	}

	if err := service.DeleteActivity(activityID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "活动删除成功",
	})
}

// SearchActivities 搜索活动
func SearchActivities(c *gin.Context) {
	keyword := c.Query("keyword")

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "搜索关键词不能为空",
		})
		return
	}

	activities, err := service.SearchActivities(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    activities,
	})
}

// GetActivityDetail 获取活动详情
func GetActivityDetail(c *gin.Context) {
	idStr := c.Param("id")
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "活动ID格式不正确",
		})
		return
	}

	activity, err := service.GetActivityDetail(activityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    activity,
	})
}
