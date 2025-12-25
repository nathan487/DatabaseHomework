package service

import (
	"log"
	"time"

	"volunteer-system/config"
	"volunteer-system/model"
)

// CloseExpiredActivities 定时任务：检查并关闭过期的活动
func CloseExpiredActivities() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		// 查询所有状态为active的活动
		var activities []model.Activity
		if err := config.DB.Where("status = ?", "active").Find(&activities).Error; err != nil {
			log.Printf("查询活动失败: %v", err)
			continue
		}

		// 逐个检查活动是否过期
		for _, activity := range activities {
			if activity.ActivityTime.Before(now) {
				// 更新状态为expired
				if err := config.DB.Model(&model.Activity{}).
					Where("activity_id = ?", activity.ActivityID).
					Update("status", "expired").Error; err != nil {
					log.Printf("更新活动状态失败 (ID:%d): %v", activity.ActivityID, err)
				} else {
					log.Printf("活动已过期，自动关闭 (ID:%d, 标题:%s)", activity.ActivityID, activity.Title)
				}
			}
		}
	}
}

// GetActiveActivities 获取所有活跃的活动（用户视角，不显示已过期的）
func GetActiveActivities(deptID, categoryID *int) ([]model.Activity, error) {
	var activities []model.Activity
	query := config.DB.Model(&model.Activity{})

	// 只查询状态为active的活动
	query = query.Where("status = ?", "active")

	if deptID != nil {
		query = query.Where("dept_id = ?", *deptID)
	}
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if err := query.Order("activity_time desc").Find(&activities).Error; err != nil {
		return nil, err
	}

	return activities, nil
}

// IsActivityExpired 检查活动是否已过期
func IsActivityExpired(activityID int) (bool, error) {
	var activity model.Activity
	if err := config.DB.First(&activity, "activity_id = ?", activityID).Error; err != nil {
		return false, err
	}

	// 如果活动时间已过或状态为expired，则认为已过期
	return activity.ActivityTime.Before(time.Now()) || activity.Status == "expired", nil
}
