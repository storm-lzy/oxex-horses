package repository

import (
	"context"

	"niuma-house/internal/model"
	"niuma-house/pkg/database"

	"gorm.io/gorm"
)

// Searcher 搜索接口
type Searcher interface {
	SearchCompanies(ctx context.Context, keyword string, page, size int) ([]model.Company, int64, error)
}

// CompanyRepository 公司仓储
type CompanyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository 创建公司仓储
func NewCompanyRepository() *CompanyRepository {
	return &CompanyRepository{db: database.GetDB()}
}

// Create 创建公司
func (r *CompanyRepository) Create(company *model.Company) error {
	return r.db.Create(company).Error
}

// FindByID 根据 ID 查找公司
func (r *CompanyRepository) FindByID(id uint) (*model.Company, error) {
	var company model.Company
	err := r.db.Preload("Creator").
		Where("status > 0").First(&company, id).Error
	if err != nil {
		return nil, err
	}
	return &company, nil
}

// Update 更新公司
func (r *CompanyRepository) Update(company *model.Company) error {
	return r.db.Save(company).Error
}

// Delete 删除公司（软删除）
func (r *CompanyRepository) Delete(id uint) error {
	return r.db.Model(&model.Company{}).Where("id = ?", id).
		Update("status", 0).Error
}

// List 公司列表
func (r *CompanyRepository) List(page, size int) ([]model.Company, int64, error) {
	var companies []model.Company
	var total int64

	r.db.Model(&model.Company{}).Where("status > 0").Count(&total)

	offset := (page - 1) * size
	err := r.db.Preload("Creator").
		Where("status > 0").
		Order("risk_level DESC, created_at DESC").
		Offset(offset).Limit(size).
		Find(&companies).Error

	return companies, total, err
}

// IncrementViews 增加浏览数
func (r *CompanyRepository) IncrementViews(companyID uint) error {
	return r.db.Model(&model.Company{}).Where("id = ?", companyID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// AdminList 管理端列表
func (r *CompanyRepository) AdminList(page, size int) ([]model.Company, int64, error) {
	var companies []model.Company
	var total int64

	r.db.Model(&model.Company{}).Count(&total)

	offset := (page - 1) * size
	err := r.db.Preload("Creator").
		Order("created_at DESC").
		Offset(offset).Limit(size).
		Find(&companies).Error

	return companies, total, err
}

// MySQLSearcher MySQL 搜索实现
type MySQLSearcher struct {
	db *gorm.DB
}

// NewMySQLSearcher 创建 MySQL 搜索器
func NewMySQLSearcher() *MySQLSearcher {
	return &MySQLSearcher{db: database.GetDB()}
}

// SearchCompanies 搜索公司
func (s *MySQLSearcher) SearchCompanies(ctx context.Context, keyword string, page, size int) ([]model.Company, int64, error) {
	var companies []model.Company
	var total int64

	query := s.db.Model(&model.Company{}).Where("status > 0")
	if keyword != "" {
		query = query.Where("name LIKE ? OR city LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * size
	err := query.Preload("Creator").
		Order("risk_level DESC, view_count DESC").
		Offset(offset).Limit(size).
		Find(&companies).Error

	return companies, total, err
}

// ElasticSearcher Elasticsearch 搜索实现（预留）
// type ElasticSearcher struct {
//     client *elastic.Client
// }
//
// func NewElasticSearcher(client *elastic.Client) *ElasticSearcher {
//     return &ElasticSearcher{client: client}
// }
//
// func (s *ElasticSearcher) SearchCompanies(ctx context.Context, keyword string, page, size int) ([]model.Company, int64, error) {
//     // TODO: 实现 Elasticsearch 搜索
//     return nil, 0, nil
// }
