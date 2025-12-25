package service

import (
	"errors"
	"strings"
	"time"

	"volunteer-system/config"
	"volunteer-system/model"

	"gorm.io/gorm"
)

func ApplyActivity(userID, activityID int) (*model.Application, error) {
	var user model.User
	if err := config.DB.First(&user, "user_id = ?", userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	var activity model.Activity
	if err := config.DB.First(&activity, "activity_id = ?", activityID).Error; err != nil {
		return nil, errors.New("活动不存在")
	}

	// 检查活动是否已过期
	if activity.Status != "active" {
		return nil, errors.New("活动已关闭，不能申请")
	}

	// 检查活动时间是否已过期
	if activity.ActivityTime.Before(time.Now()) {
		return nil, errors.New("活动已过期，不能申请")
	}

	var existingCount int64
	if err := config.DB.Model(&model.Application{}).
		Where("user_id = ? AND activity_id = ?", userID, activityID).
		Count(&existingCount).Error; err != nil {
		return nil, errors.New("查询报名记录失败")
	}
	if existingCount > 0 {
		return nil, errors.New("您已申请参加该活动")
	}

	var approvedCount int64
	if err := config.DB.Model(&model.Application{}).
		Where("activity_id = ? AND current_status = ?", activityID, "approved").
		Count(&approvedCount).Error; err != nil {
		return nil, errors.New("查询活动报名人数失败")
	}
	if approvedCount >= int64(activity.MaxPeople) {
		return nil, errors.New("活动人数已满")
	}

	now := time.Now()

	application := model.Application{
		UserID:        userID,
		ActivityID:    activityID,
		ApplyTime:     now,
		CurrentStatus: "pending",
	}

	if err := config.DB.Create(&application).Error; err != nil {
		return nil, errors.New("报名失败")
	}

	log := model.ApplicationStatusLog{
		ApplicationID: application.ApplicationID,
		HandlerID:     &userID,
		LogStatus:     "pending",
		HandleTime:    now,
	}

	if err := config.DB.Create(&log).Error; err != nil {
		return nil, errors.New("保存报名日志失败")
	}

	return &application, nil
}

func ListActivityApplications(activityID int) ([]model.ActivityApplicationWithUser, error) {
	var apps []model.ActivityApplicationWithUser
	if err := config.DB.Table("Application").
		Select("Application.application_id, Application.user_id, User.username, Application.apply_time, Application.current_status").
		Joins("JOIN User ON Application.user_id = User.user_id").
		Where("Application.activity_id = ?", activityID).
		Order("Application.apply_time DESC").
		Scan(&apps).Error; err != nil {
		return nil, errors.New("查询报名记录失败")
	}

	return apps, nil
}

func ListUserApplications(userID int) ([]model.UserApplicationInfo, error) {
	var apps []model.UserApplicationInfo
	if err := config.DB.Table("Application").
		Select("Application.application_id, Application.activity_id, Activity.title, Activity.activity_time, Activity.location, Application.current_status, Application.apply_time").
		Joins("JOIN Activity ON Application.activity_id = Activity.activity_id").
		Where("Application.user_id = ?", userID).
		Order("Application.apply_time DESC").
		Scan(&apps).Error; err != nil {
		return nil, errors.New("查询报名记录失败")
	}

	return apps, nil
}

func UpdateApplicationStatus(appID int, status string, handlerID int) error {
	status = strings.ToLower(strings.TrimSpace(status))
	if status != "approved" && status != "rejected" && status != "pending" {
		return errors.New("状态只能是 approved / rejected / pending")
	}

	var app model.Application
	if err := config.DB.First(&app, "application_id = ?", appID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("报名记录不存在")
		}
		return errors.New("查询报名记录失败")
	}

	if status == "approved" && app.CurrentStatus != "approved" {
		var activity model.Activity
		if err := config.DB.First(&activity, "activity_id = ?", app.ActivityID).Error; err != nil {
			return errors.New("查询活动信息失败")
		}

		var approvedCount int64
		if err := config.DB.Model(&model.Application{}).
			Where("activity_id = ? AND current_status = ?", app.ActivityID, "approved").
			Count(&approvedCount).Error; err != nil {
			return errors.New("查询活动报名人数失败")
		}

		if approvedCount >= int64(activity.MaxPeople) {
			return errors.New("活动人数已满，无法再通过报名")
		}
	}

	now := time.Now()

	if err := config.DB.Model(&model.Application{}).
		Where("application_id = ?", appID).
		Update("current_status", status).Error; err != nil {
		return errors.New("更新报名状态失败")
	}

	log := model.ApplicationStatusLog{
		ApplicationID: appID,
		HandlerID:     &handlerID,
		LogStatus:     status,
		HandleTime:    now,
	}

	if err := config.DB.Create(&log).Error; err != nil {
		return errors.New("保存审核日志失败")
	}

	return nil
}

// CancelApplication 取消报名
func CancelApplication(appID int) error {
	var app model.Application
	if err := config.DB.First(&app, "application_id = ?", appID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("报名记录不存在")
		}
		return errors.New("查询报名记录失败")
	}

	// 检查活动是否已开始
	var activity model.Activity
	if err := config.DB.First(&activity, "activity_id = ?", app.ActivityID).Error; err != nil {
		return errors.New("查询活动信息失败")
	}

	if activity.ActivityTime.Before(time.Now()) {
		return errors.New("活动已开始，不能取消报名")
	}

	// 只允许取消待审批和已批准的报名
	if app.CurrentStatus != "pending" && app.CurrentStatus != "approved" {
		return errors.New("该报名状态不允许取消")
	}

	// 先删除相关的状态日志（因为有外键约束）
	if err := config.DB.Delete(&model.ApplicationStatusLog{}, "application_id = ?", appID).Error; err != nil {
		return errors.New("删除状态日志失败")
	}

	// 再删除报名记录
	if err := config.DB.Delete(&model.Application{}, "application_id = ?", appID).Error; err != nil {
		return errors.New("取消报名失败")
	}

	return nil
}
