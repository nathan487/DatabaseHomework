package service

import (
	"errors"

	"volunteer-system/config"
	"volunteer-system/model"
	"volunteer-system/utils"

	"gorm.io/gorm"
)

func ListActivities(deptID, categoryID *int) ([]model.Activity, error) {
	var activities []model.Activity
	query := config.DB.Model(&model.Activity{})

	// 只返回活跃的活动，不显示已过期或已关闭的活动
	query = query.Where("status = ?", "active")

	if deptID != nil {
		query = query.Where("dept_id = ?", *deptID)
	}
	if categoryID != nil {
		query = query.Where("category_id = ?", *categoryID)
	}

	if err := query.Order("activity_time desc").Find(&activities).Error; err != nil {
		return nil, errors.New("查询活动失败")
	}

	return activities, nil
}

func CreateActivity(req *model.CreateActivityRequest) (*model.Activity, error) {
	activityTime, err := utils.ParseActivityTime(req.ActivityTime)
	if err != nil {
		return nil, errors.New("活动时间格式不正确")
	}

	activity := model.Activity{
		DeptID:       req.DeptID,
		CategoryID:   req.CategoryID,
		CreatorID:    req.CreatorID,
		Title:        req.Title,
		Description:  req.Description,
		ActivityTime: activityTime,
		Location:     req.Location,
		MaxPeople:    req.MaxPeople,
	}

	if err := config.DB.Create(&activity).Error; err != nil {
		return nil, errors.New("创建活动失败")
	}

	return &activity, nil
}

func UpdateActivity(activityID int, req *model.UpdateActivityRequest) (*model.Activity, error) {
	activityTime, err := utils.ParseActivityTime(req.ActivityTime)
	if err != nil {
		return nil, errors.New("活动时间格式不正确")
	}

	var activity model.Activity
	if err := config.DB.First(&activity, "activity_id = ?", activityID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("活动不存在")
		}
		return nil, errors.New("查询活动失败")
	}

	activity.DeptID = req.DeptID
	activity.CategoryID = req.CategoryID
	activity.CreatorID = req.CreatorID
	activity.Title = req.Title
	activity.Description = req.Description
	activity.ActivityTime = activityTime
	activity.Location = req.Location
	activity.MaxPeople = req.MaxPeople

	if err := config.DB.Save(&activity).Error; err != nil {
		return nil, errors.New("更新活动失败")
	}

	return &activity, nil
}

func DeleteActivity(activityID int) error {
	// 先删除所有相关的应用状态日志
	if err := config.DB.Delete(&model.ApplicationStatusLog{}, "application_id IN (SELECT application_id FROM Application WHERE activity_id = ?)", activityID).Error; err != nil {
		return errors.New("删除状态日志失败")
	}

	// 再删除所有相关的报名记录
	if err := config.DB.Delete(&model.Application{}, "activity_id = ?", activityID).Error; err != nil {
		return errors.New("删除报名记录失败")
	}

	// 最后删除活动
	if err := config.DB.Delete(&model.Activity{}, "activity_id = ?", activityID).Error; err != nil {
		return errors.New("删除活动失败")
	}
	return nil
}

// SearchActivities 搜索活动
func SearchActivities(keyword string) ([]model.Activity, error) {
	var activities []model.Activity
	if err := config.DB.Where("status = ? AND title LIKE ?", "active", "%"+keyword+"%").
		Order("activity_time desc").
		Find(&activities).Error; err != nil {
		return nil, errors.New("搜索活动失败")
	}
	return activities, nil
}

// GetActivityDetail 获取活动详情
func GetActivityDetail(activityID int) (*model.Activity, error) {
	var activity model.Activity
	if err := config.DB.First(&activity, "activity_id = ?", activityID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("活动不存在")
		}
		return nil, errors.New("获取活动详情失败")
	}
	return &activity, nil
}
