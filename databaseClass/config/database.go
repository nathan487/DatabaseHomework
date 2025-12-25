package config

import (
	"errors"

	"volunteer-system/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	dsn := "root:123456@tcp(127.0.0.1:3306)/volunteer?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return ensureDefaultRoles()
}

func ensureDefaultRoles() error {
	defaultRoleNames := []string{"admin", "user"}

	for _, roleName := range defaultRoleNames {
		var existing model.Role
		err := DB.Where("role_name = ?", roleName).First(&existing).Error
		if err == nil {
			continue
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		role := model.Role{
			RoleName: roleName,
		}
		if err := DB.Create(&role).Error; err != nil {
			return err
		}
	}

	return nil
}
