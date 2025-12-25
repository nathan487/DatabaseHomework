package service

import (
	"errors"
	"strings"

	"volunteer-system/config"
	"volunteer-system/model"
	"volunteer-system/utils"

	"gorm.io/gorm"
)

func Register(username, password, roleName string) (*model.LoginResponse, error) {
	var existing model.User
	if err := config.DB.Where("username = ?", username).First(&existing).Error; err == nil {
		return nil, errors.New("用户名已存在")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("查询用户失败")
	}

	if roleName == "" {
		roleName = "user"
	} else {
		roleName = strings.TrimSpace(roleName)
	}

	var role model.Role
	if err := config.DB.Where("role_name = ?", roleName).First(&role).Error; err != nil {
		return nil, errors.New("角色不存在")
	}

	user := model.User{
		RoleID:   role.RoleID,
		Username: username,
		Password: utils.MD5Hash(password),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, errors.New("注册失败")
	}

	return &model.LoginResponse{
		UserID:   user.UserID,
		Username: user.Username,
		RoleName: role.RoleName,
	}, nil
}

func Login(username, password string) (*model.LoginResponse, error) {
	var user model.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Password != utils.MD5Hash(password) {
		return nil, errors.New("用户名或密码错误")
	}

	var role model.Role
	if err := config.DB.First(&role, "role_id = ?", user.RoleID).Error; err != nil {
		return nil, errors.New("查询角色失败")
	}

	return &model.LoginResponse{
		UserID:   user.UserID,
		Username: user.Username,
		RoleName: role.RoleName,
	}, nil
}
