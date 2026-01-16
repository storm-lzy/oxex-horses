package repository

import (
	"niuma-house/internal/model"
	"niuma-house/pkg/database"
)

// OccupationRepository 职业仓储
type OccupationRepository struct{}

// NewOccupationRepository 创建职业仓储
func NewOccupationRepository() *OccupationRepository {
	return &OccupationRepository{}
}

// List 获取所有职业
func (r *OccupationRepository) List() ([]model.Occupation, error) {
	var occupations []model.Occupation
	err := database.GetDB().Order("id ASC").Find(&occupations).Error
	return occupations, err
}

// FindByID 根据 ID 查找职业
func (r *OccupationRepository) FindByID(id uint) (*model.Occupation, error) {
	var occupation model.Occupation
	err := database.GetDB().First(&occupation, id).Error
	if err != nil {
		return nil, err
	}
	return &occupation, nil
}
