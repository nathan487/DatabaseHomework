package model

import "time"

type Role struct {
	RoleID   int    `json:"role_id" gorm:"column:role_id;primaryKey;autoIncrement"`
	RoleName string `json:"role_name" gorm:"column:role_name;not null"`
}

func (Role) TableName() string {
	return "Role"
}

type User struct {
	UserID   int    `json:"user_id" gorm:"column:user_id;primaryKey;autoIncrement"`
	RoleID   int    `json:"role_id" gorm:"column:role_id;not null"`
	Username string `json:"username" gorm:"column:username;not null;unique"`
	Password string `json:"-" gorm:"column:password;not null"`
}

func (User) TableName() string {
	return "User"
}

type Dept struct {
	DeptID   int    `json:"dept_id" gorm:"column:dept_id;primaryKey;autoIncrement"`
	DeptName string `json:"dept_name" gorm:"column:dept_name;not null"`
}

func (Dept) TableName() string {
	return "Dept"
}

type ActivityCategory struct {
	CategoryID   int    `json:"category_id" gorm:"column:category_id;primaryKey;autoIncrement"`
	CategoryName string `json:"category_name" gorm:"column:category_name;not null"`
}

func (ActivityCategory) TableName() string {
	return "ActivityCategory"
}

type Activity struct {
	ActivityID   int       `json:"activity_id" gorm:"column:activity_id;primaryKey;autoIncrement"`
	DeptID       int       `json:"dept_id" gorm:"column:dept_id;not null"`
	CategoryID   int       `json:"category_id" gorm:"column:category_id;not null"`
	CreatorID    int       `json:"creator_id" gorm:"column:creator_id;not null"`
	Title        string    `json:"title" gorm:"column:title;not null"`
	Description  string    `json:"description" gorm:"column:description;type:longtext"`
	ActivityTime time.Time `json:"activity_time" gorm:"column:activity_time;not null"`
	Location     string    `json:"location" gorm:"column:location;not null"`
	MaxPeople    int       `json:"max_people" gorm:"column:max_people;not null"`
	Status       string    `json:"status" gorm:"column:status;default:active"`
}

func (Activity) TableName() string {
	return "Activity"
}

type Application struct {
	ApplicationID int       `json:"application_id" gorm:"column:application_id;primaryKey;autoIncrement"`
	UserID        int       `json:"user_id" gorm:"column:user_id;not null"`
	ActivityID    int       `json:"activity_id" gorm:"column:activity_id;not null"`
	ApplyTime     time.Time `json:"apply_time" gorm:"column:apply_time;not null"`
	CurrentStatus string    `json:"current_status" gorm:"column:current_status;not null;default:pending"`
}

func (Application) TableName() string {
	return "Application"
}

type ApplicationStatusLog struct {
	LogID         int       `json:"log_id" gorm:"column:log_id;primaryKey;autoIncrement"`
	ApplicationID int       `json:"application_id" gorm:"column:application_id;not null"`
	HandlerID     *int      `json:"handler_id" gorm:"column:handler_id"`
	LogStatus     string    `json:"log_status" gorm:"column:log_status;not null"`
	HandleTime    time.Time `json:"handle_time" gorm:"column:handle_time;not null"`
}

func (ApplicationStatusLog) TableName() string {
	return "ApplicationStatusLog"
}

// Request and Response structs
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RoleName string `json:"role_name"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	RoleName string `json:"role_name"`
	RoleID   int    `json:"role_id"`
}

type CreateActivityRequest struct {
	DeptID       int    `json:"dept_id" binding:"required"`
	CategoryID   int    `json:"category_id" binding:"required"`
	CreatorID    int    `json:"creator_id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	ActivityTime string `json:"activity_time" binding:"required"`
	Location     string `json:"location" binding:"required"`
	MaxPeople    int    `json:"max_people" binding:"required"`
}

type UpdateActivityRequest struct {
	DeptID       int    `json:"dept_id" binding:"required"`
	CategoryID   int    `json:"category_id" binding:"required"`
	CreatorID    int    `json:"creator_id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	ActivityTime string `json:"activity_time" binding:"required"`
	Location     string `json:"location" binding:"required"`
	MaxPeople    int    `json:"max_people" binding:"required"`
}

type ApplyActivityRequest struct {
	UserID int `json:"user_id" binding:"required"`
}

type UpdateApplicationStatusRequest struct {
	Status    string `json:"status" binding:"required"`
	HandlerID int    `json:"handler_id" binding:"required"`
}

type ActivityApplicationWithUser struct {
	ApplicationID int       `json:"application_id"`
	UserID        int       `json:"user_id"`
	Username      string    `json:"username"`
	ApplyTime     time.Time `json:"apply_time"`
	CurrentStatus string    `json:"current_status"`
}

type UserApplicationInfo struct {
	ApplicationID int       `json:"application_id"`
	ActivityID    int       `json:"activity_id"`
	Title         string    `json:"title"`
	ActivityTime  time.Time `json:"activity_time"`
	Location      string    `json:"location"`
	CurrentStatus string    `json:"current_status"`
	ApplyTime     time.Time `json:"apply_time"`
}
