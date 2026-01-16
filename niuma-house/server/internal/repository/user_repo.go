package repository

import (
	"niuma-house/internal/model"
	"niuma-house/pkg/database"

	"gorm.io/gorm"
)

// UserRepository 用户仓储
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.GetDB()}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据 ID 查找用户
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Occupation").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// UpdateExp 更新经验值
func (r *UserRepository) UpdateExp(userID uint, expDelta int) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		UpdateColumn("exp", gorm.Expr("exp + ?", expDelta)).Error
}

// UpdateLevel 更新等级
func (r *UserRepository) UpdateLevel(userID uint, level int) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("level", level).Error
}

// Ban 封禁用户
func (r *UserRepository) Ban(userID uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("status", 0).Error
}

// Unban 解封用户
func (r *UserRepository) Unban(userID uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("status", 1).Error
}

// List 用户列表
func (r *UserRepository) List(page, size int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	r.db.Model(&model.User{}).Count(&total)

	offset := (page - 1) * size
	err := r.db.Preload("Occupation").
		Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&users).Error

	return users, total, err
}
