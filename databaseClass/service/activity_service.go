package service

import (
	"errors"

	"volunteer-system/config"
	"volunteer-system/model"
	"volunteer-system/utils"

	"gorm.io/gorm"
)

// ActivityDetail 活动详情（4表JOIN）
type ActivityDetail struct {
	ActivityID   int    `json:"activity_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Location     string `json:"location"`
	ActivityTime string `json:"activity_time"`
	MaxPeople    int    `json:"max_people"`
	Status       string `json:"status"`
	DeptID       int    `json:"dept_id"`
	DeptName     string `json:"dept_name"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	CreatorID    int    `json:"creator_id"`
	CreatorName  string `json:"creator_name"`
}

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

// GetActivityDetail 获取活动详情（4表JOIN：活动+部门+分类+创建者）
func GetActivityDetail(activityID int) (*ActivityDetail, error) {
	var activity ActivityDetail

	err := config.DB.Raw(`
		SELECT a.activity_id, a.title, a.description, a.location,
			DATE_FORMAT(a.activity_time, '%Y-%m-%d %H:%i') as activity_time,
			a.max_people, a.status,
			a.dept_id, COALESCE(d.dept_name, '') as dept_name,
			a.category_id, COALESCE(ac.category_name, '') as category_name,
			a.creator_id, COALESCE(u.username, '') as creator_name
		FROM Activity a
		LEFT JOIN Dept d ON a.dept_id = d.dept_id
		LEFT JOIN ActivityCategory ac ON a.category_id = ac.category_id
		LEFT JOIN User u ON a.creator_id = u.user_id
		WHERE a.activity_id = ?
	`, activityID).Scan(&activity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("活动不存在")
		}
		return nil, errors.New("获取活动详情失败")
	}

	// 检查是否有结果
	if activity.ActivityID == 0 {
		return nil, errors.New("活动不存在")
	}

	return &activity, nil
}

// AvailableActivity 用户可申请的活动（NOT IN集合操作）
type AvailableActivity struct {
	ActivityID        int    `json:"activity_id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Location          string `json:"location"`
	ActivityTime      string `json:"activity_time"`
	MaxPeople         int    `json:"max_people"`
	CurrentApplyCount int    `json:"current_apply_count"`
	RemainingSlots    int    `json:"remaining_slots"`
	DeptName          string `json:"dept_name"`
	CategoryName      string `json:"category_name"`
}

// GetAvailableActivities 获取用户可申请的活动（NOT IN集合操作+自连接：有空位+未申请+无时间冲突）
func GetAvailableActivities(userID int) ([]AvailableActivity, error) {
	var activities []AvailableActivity

	err := config.DB.Raw(`
		SELECT a.activity_id, a.title, a.description, a.location,
			DATE_FORMAT(a.activity_time, '%Y-%m-%d %H:%i') as activity_time,
			a.max_people, COALESCE(COUNT(app.application_id), 0) as current_apply_count,
			(a.max_people - COALESCE(COUNT(app.application_id), 0)) as remaining_slots,
			COALESCE(d.dept_name, '未分配') as dept_name,
			COALESCE(ac.category_name, '未分类') as category_name
		FROM Activity a
		LEFT JOIN Application app ON a.activity_id = app.activity_id
		LEFT JOIN Dept d ON a.dept_id = d.dept_id
		LEFT JOIN ActivityCategory ac ON a.category_id = ac.category_id
		WHERE a.status = 'active'
		AND a.activity_id NOT IN (
			SELECT DISTINCT activity_id FROM Application WHERE user_id = ?
		)
		AND a.activity_id NOT IN (
			-- 自连接：排除与用户已申请活动时间冲突的
			SELECT DISTINCT a2.activity_id
			FROM Activity a2
			JOIN Activity a1 ON ABS(HOUR(TIMEDIFF(a1.activity_time, a2.activity_time))) < 2
			WHERE a1.activity_id IN (
				SELECT DISTINCT activity_id FROM Application WHERE user_id = ?
			) AND a1.status = 'active' AND a2.status = 'active'
		)
		GROUP BY a.activity_id, a.title, a.description, a.location, a.activity_time, 
			a.max_people, d.dept_name, ac.category_name
		HAVING remaining_slots > 0
		ORDER BY a.activity_time ASC
	`, userID, userID).Scan(&activities).Error

	if err != nil {
		return nil, errors.New("获取可申请活动列表失败")
	}

	return activities, nil
}
