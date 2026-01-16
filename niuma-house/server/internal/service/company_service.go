package service

import (
	"context"

	"niuma-house/internal/model"
	"niuma-house/internal/repository"
)

// CompanyService 公司服务
type CompanyService struct {
	companyRepo *repository.CompanyRepository
	searcher    repository.Searcher
}

// NewCompanyService 创建公司服务
func NewCompanyService() *CompanyService {
	return &CompanyService{
		companyRepo: repository.NewCompanyRepository(),
		searcher:    repository.NewMySQLSearcher(), // 使用 MySQL 搜索实现
	}
}

// CreateCompanyRequest 创建公司请求
type CreateCompanyRequest struct {
	Name      string   `json:"name" binding:"required,max=100"`
	City      string   `json:"city" binding:"max=50"`
	Tags      []string `json:"tags"`
	RiskLevel int      `json:"risk_level" binding:"min=1,max=5"`
	Evidence  []string `json:"evidence"`
	Content   string   `json:"content"`
}

// Create 创建公司
func (s *CompanyService) Create(userID uint, req *CreateCompanyRequest) (*model.Company, error) {
	company := &model.Company{
		Name:      req.Name,
		City:      req.City,
		Tags:      req.Tags,
		RiskLevel: req.RiskLevel,
		Evidence:  req.Evidence,
		Content:   req.Content,
		CreatorID: userID,
		Status:    1,
	}

	if err := s.companyRepo.Create(company); err != nil {
		return nil, err
	}

	return company, nil
}

// GetByID 获取公司详情
func (s *CompanyService) GetByID(id uint) (*model.Company, error) {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 增加浏览量
	s.companyRepo.IncrementViews(id)

	return company, nil
}

// List 公司列表
func (s *CompanyService) List(page, size int) ([]model.Company, int64, error) {
	return s.companyRepo.List(page, size)
}

// Search 搜索公司
func (s *CompanyService) Search(ctx context.Context, keyword string, page, size int) ([]model.Company, int64, error) {
	return s.searcher.SearchCompanies(ctx, keyword, page, size)
}

// AdminList 管理端列表
func (s *CompanyService) AdminList(page, size int) ([]model.Company, int64, error) {
	return s.companyRepo.AdminList(page, size)
}

// AdminDelete 管理员删除
func (s *CompanyService) AdminDelete(companyID uint) error {
	return s.companyRepo.Delete(companyID)
}
