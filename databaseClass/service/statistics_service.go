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

// 按部门统计
type DeptStatistics struct {
	DeptID        int     `json:"dept_id"`
	DeptName      string  `json:"dept_name"`
	ActivityCount int     `json:"activity_count"`
	AvgCapacity   float64 `json:"avg_capacity"`
	CreatorCount  int     `json:"creator_count"`
}

// 按分类统计
type CategoryStatistics struct {
	CategoryID        int    `json:"category_id"`
	CategoryName      string `json:"category_name"`
	ActivityCount     int    `json:"activity_count"`
	TotalApplications int    `json:"total_applications"`
	ApprovedCount     int    `json:"approved_count"`
}

// 用户活跃度统计
type UserActivityStatistics struct {
	UserID        int    `json:"user_id"`
	Username      string `json:"username"`
	TotalApplied  int    `json:"total_applied"`
	ApprovedCount int    `json:"approved_count"`
	RejectedCount int    `json:"rejected_count"`
}

// 活动热度排行
type ActivityPopularity struct {
	ActivityID       int     `json:"activity_id"`
	Title            string  `json:"title"`
	MaxPeople        int     `json:"max_people"`
	ApplicationCount int     `json:"application_count"`
	ApprovedCount    int     `json:"approved_count"`
	FillRate         float64 `json:"fill_rate"`
}

// 管理员创建统计
type AdminCreationStatistics struct {
	UserID               int    `json:"user_id"`
	Username             string `json:"username"`
	CreatedActivities    int    `json:"created_activities"`
	TotalApplications    int    `json:"total_applications"`
	ApprovedApplications int    `json:"approved_applications"`
}

// 即将开始的活动（日期时间函数）
type UpcomingActivity struct {
	ActivityID     int    `json:"activity_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	DeptName       string `json:"dept_name"`
	CategoryName   string `json:"category_name"`
	MaxPeople      int    `json:"max_people"`
	ActivityTime   string `json:"activity_time"`
	DaysRemaining  int    `json:"days_remaining"`
	HoursRemaining int    `json:"hours_remaining"`
}

// 热门活动（关联子查询）
type PopularActivity struct {
	ActivityID    int     `json:"activity_id"`
	Title         string  `json:"title"`
	DeptName      string  `json:"dept_name"`
	CategoryName  string  `json:"category_name"`
	CreatorName   string  `json:"creator_name"`
	MaxPeople     int     `json:"max_people"`
	ApprovedCount int     `json:"approved_count"`
	PendingCount  int     `json:"pending_count"`
	RejectedCount int     `json:"rejected_count"`
	FillRate      float64 `json:"fill_rate"`
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

// GetDeptStatistics 获取按部门统计的数据 (聚合函数+GROUP BY)
func GetDeptStatistics() ([]DeptStatistics, error) {
	var results []DeptStatistics

	err := config.DB.Raw(`
		SELECT d.dept_id, d.dept_name,
			COUNT(DISTINCT a.activity_id) as activity_count,
			IFNULL(AVG(a.max_people), 0) as avg_capacity,
			COUNT(DISTINCT a.creator_id) as creator_count
		FROM Dept d
		LEFT JOIN Activity a ON d.dept_id = a.dept_id
		GROUP BY d.dept_id, d.dept_name
		ORDER BY activity_count DESC
	`).Scan(&results).Error

	if err != nil {
		return nil, errors.New("查询部门统计失败")
	}

	return results, nil
}

// GetCategoryStatistics 获取按分类统计的数据 (聚合函数+GROUP BY)
func GetCategoryStatistics() ([]CategoryStatistics, error) {
	var results []CategoryStatistics

	err := config.DB.Raw(`
		SELECT ac.category_id, ac.category_name,
			COUNT(DISTINCT a.activity_id) as activity_count,
			COALESCE(COUNT(ap.application_id), 0) as total_applications,
			COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) as approved_count
		FROM ActivityCategory ac
		LEFT JOIN Activity a ON ac.category_id = a.category_id
		LEFT JOIN Application ap ON a.activity_id = ap.activity_id
		GROUP BY ac.category_id, ac.category_name
		ORDER BY activity_count DESC
	`).Scan(&results).Error

	if err != nil {
		return nil, errors.New("查询分类统计失败")
	}

	return results, nil
}

// GetUserActivityStatistics 获取用户活跃度统计 (聚合函数+GROUP BY+HAVING)
func GetUserActivityStatistics() ([]UserActivityStatistics, error) {
	var results []UserActivityStatistics

	err := config.DB.Raw(`
		SELECT u.user_id, u.username,
			COUNT(ap.application_id) as total_applied,
			COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) as approved_count,
			COALESCE(SUM(CASE WHEN ap.current_status = 'rejected' THEN 1 ELSE 0 END), 0) as rejected_count
		FROM User u
		LEFT JOIN Application ap ON u.user_id = ap.user_id
		GROUP BY u.user_id, u.username
		HAVING COUNT(ap.application_id) > 0
		ORDER BY approved_count DESC
	`).Scan(&results).Error

	if err != nil {
		return nil, errors.New("查询用户活跃度统计失败")
	}

	return results, nil
}

// GetActivityPopularity 获取活动热度排行 (聚合函数+GROUP BY)
func GetActivityPopularity() ([]ActivityPopularity, error) {
	var results []ActivityPopularity

	err := config.DB.Raw(`
		SELECT a.activity_id, a.title, a.max_people,
			COALESCE(COUNT(ap.application_id), 0) as application_count,
			COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) as approved_count,
			ROUND(COALESCE(COUNT(ap.application_id), 0) / a.max_people * 100, 2) as fill_rate
		FROM Activity a
		LEFT JOIN Application ap ON a.activity_id = ap.activity_id
		GROUP BY a.activity_id, a.title, a.max_people
		ORDER BY application_count DESC
		LIMIT 10
	`).Scan(&results).Error

	if err != nil {
		return nil, errors.New("查询活动热度排行失败")
	}

	return results, nil
}

// GetAdminCreationStatistics 获取管理员创建活动统计 (聚合函数+GROUP BY)
func GetAdminCreationStatistics() ([]AdminCreationStatistics, error) {
	var results []AdminCreationStatistics

	err := config.DB.Raw(`
		SELECT u.user_id, u.username,
			COUNT(DISTINCT a.activity_id) as created_activities,
			COALESCE(COUNT(ap.application_id), 0) as total_applications,
			COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) as approved_applications
		FROM User u
		LEFT JOIN Activity a ON u.user_id = a.creator_id
		LEFT JOIN Application ap ON a.activity_id = ap.activity_id
		WHERE u.role_id = 2
		GROUP BY u.user_id, u.username
		ORDER BY created_activities DESC
	`).Scan(&results).Error

	if err != nil {
		return nil, errors.New("查询管理员创建统计失败")
	}

	return results, nil
}

// GetUpcomingActivities 获取即将开始的活动 (日期时间函数: NOW(), DATE_ADD(), DATEDIFF())
func GetUpcomingActivities() ([]UpcomingActivity, error) {
	var results []UpcomingActivity

	err := config.DB.Raw(`
		SELECT a.activity_id, a.title, a.description, 
			d.dept_name, ac.category_name, a.max_people,
			DATE_FORMAT(a.activity_time, '%Y-%m-%d %H:%i') as activity_time,
			DATEDIFF(a.activity_time, NOW()) as days_remaining,
			HOUR(TIMEDIFF(a.activity_time, NOW())) as hours_remaining
		FROM Activity a
		LEFT JOIN Dept d ON a.dept_id = d.dept_id
		LEFT JOIN ActivityCategory ac ON a.category_id = ac.category_id
		WHERE a.activity_time > NOW()
			AND a.activity_time <= DATE_ADD(NOW(), INTERVAL 3 DAY)
			AND a.status = 'active'
		ORDER BY a.activity_time ASC
	`).Scan(&results).Error

	if err != nil {
		return nil, errors.New("查询即将开始的活动失败")
	}

	return results, nil
}

// GetPopularActivities 获取热门活动（关联子查询：EXISTS）
func GetPopularActivities() ([]PopularActivity, error) {
	var results []PopularActivity

	err := config.DB.Raw(`
		SELECT a.activity_id, a.title, 
			COALESCE(d.dept_name, '') as dept_name, 
			COALESCE(ac.category_name, '') as category_name, 
			COALESCE(u.username, '') as creator_name,
			a.max_people,
			COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) as approved_count,
			COALESCE(SUM(CASE WHEN ap.current_status = 'pending' THEN 1 ELSE 0 END), 0) as pending_count,
			COALESCE(SUM(CASE WHEN ap.current_status = 'rejected' THEN 1 ELSE 0 END), 0) as rejected_count,
			ROUND(COALESCE(SUM(CASE WHEN ap.current_status = 'approved' THEN 1 ELSE 0 END), 0) / a.max_people * 100, 2) as fill_rate
		FROM Activity a
		LEFT JOIN Dept d ON a.dept_id = d.dept_id
		LEFT JOIN ActivityCategory ac ON a.category_id = ac.category_id
		LEFT JOIN User u ON a.creator_id = u.user_id
		LEFT JOIN Application ap ON a.activity_id = ap.activity_id
		WHERE a.status = 'active'
			AND EXISTS (
				SELECT 1 FROM Application ap2
				WHERE a.activity_id = ap2.activity_id 
				AND ap2.current_status = 'approved'
			)
		GROUP BY a.activity_id, a.title, a.max_people, a.dept_id, a.category_id, a.creator_id
		ORDER BY approved_count DESC
		LIMIT 10
	`).Scan(&results).Error

	if err != nil {
		return nil, errors.New("查询热门活动失败: " + err.Error())
	}

	return results, nil
}

// OmnipotentVolunteer 全能志愿者（参加了所有分类活动）
type OmnipotentVolunteer struct {
	UserID                 int    `json:"user_id"`
	Username               string `json:"username"`
	CategoriesParticipated int    `json:"categories_participated"`
	TotalCategoriesCount   int    `json:"total_categories_count"`
	ApprovedCount          int    `json:"approved_count"`
}

// GetOmnipotentVolunteers 查询全能志愿者（除法查询：参加了所有分类活动的用户）
func GetOmnipotentVolunteers() ([]OmnipotentVolunteer, error) {
	var results []OmnipotentVolunteer

	err := config.DB.Raw(`
		SELECT 
			u.user_id, 
			u.username,
			COUNT(DISTINCT a.category_id) as categories_participated,
			(SELECT COUNT(*) FROM ActivityCategory) as total_categories_count,
			COUNT(DISTINCT app.application_id) as approved_count
		FROM User u
		JOIN Application app ON u.user_id = app.user_id AND app.current_status = 'approved'
		JOIN Activity a ON app.activity_id = a.activity_id
		GROUP BY u.user_id, u.username
		HAVING COUNT(DISTINCT a.category_id) = (SELECT COUNT(*) FROM ActivityCategory)
		ORDER BY approved_count DESC
	`).Scan(&results).Error

	if err != nil {
		return nil, errors.New("查询全能志愿者失败")
	}

	return results, nil
}
