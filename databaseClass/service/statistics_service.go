package service

import (
	"errors"

	"volunteer-system/config"
	"volunteer-system/model"
)

type StatisticsData struct {
	TotalActivities      int64 `json:"total_activities"`
	TotalUsers           int64 `json:"total_users"`
	TotalApplications    int64 `json:"total_applications"`
	ApprovedApplications int64 `json:"approved_applications"`
	PendingApplications  int64 `json:"pending_applications"`
	RejectedApplications int64 `json:"rejected_applications"`
	ActiveActivities     int64 `json:"active_activities"`
	ExpiredActivities    int64 `json:"expired_activities"`
}

// GetStatistics 获取系统统计数据
func GetStatistics() (*StatisticsData, error) {
	stats := &StatisticsData{}

	// 总活动数
	if err := config.DB.Model(&model.Activity{}).Count(&stats.TotalActivities).Error; err != nil {
		return nil, errors.New("查询总活动数失败")
	}

	// 总用户数
	if err := config.DB.Model(&model.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, errors.New("查询总用户数失败")
	}

	// 总报名数
	if err := config.DB.Model(&model.Application{}).Count(&stats.TotalApplications).Error; err != nil {
		return nil, errors.New("查询总报名数失败")
	}

	// 已批准报名数
	if err := config.DB.Model(&model.Application{}).
		Where("current_status = ?", "approved").
		Count(&stats.ApprovedApplications).Error; err != nil {
		return nil, errors.New("查询已批准报名数失败")
	}

	// 待审批报名数
	if err := config.DB.Model(&model.Application{}).
		Where("current_status = ?", "pending").
		Count(&stats.PendingApplications).Error; err != nil {
		return nil, errors.New("查询待审批报名数失败")
	}

	// 已拒绝报名数
	if err := config.DB.Model(&model.Application{}).
		Where("current_status = ?", "rejected").
		Count(&stats.RejectedApplications).Error; err != nil {
		return nil, errors.New("查询已拒绝报名数失败")
	}

	// 活跃活动数
	if err := config.DB.Model(&model.Activity{}).
		Where("status = ?", "active").
		Count(&stats.ActiveActivities).Error; err != nil {
		return nil, errors.New("查询活跃活动数失败")
	}

	// 已过期活动数
	if err := config.DB.Model(&model.Activity{}).
		Where("status = ?", "expired").
		Count(&stats.ExpiredActivities).Error; err != nil {
		return nil, errors.New("查询已过期活动数失败")
	}

	return stats, nil
}
