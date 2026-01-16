package service

import (
	"errors"

	"niuma-house/internal/model"
	"niuma-house/internal/repository"
	"niuma-house/pkg/jwt"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		userRepo: repository.NewUserRepository(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username     string `json:"username" binding:"required,min=3,max=20"`
	Password     string `json:"password" binding:"required,min=6,max=32"`
	OccupationID uint   `json:"occupation_id" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Register 用户注册
func (s *UserService) Register(req *RegisterRequest) (*model.User, error) {
	// 检查用户名是否存在
	_, err := s.userRepo.FindByUsername(req.Username)
	if err == nil {
		return nil, errors.New("用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username:     req.Username,
		Password:     req.Password,
		OccupationID: req.OccupationID,
		Level:        1,
		Exp:          0,
		Role:         "user",
		Status:       1,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 清除密码
	user.Password = ""
	return user, nil
}

// Login 用户登录
func (s *UserService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, errors.New("账号已被封禁")
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成 Token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	// 清除密码
	user.Password = ""

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetProfile 获取用户资料
func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(userID uint, occupationID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.OccupationID = occupationID
	return s.userRepo.Update(user)
}

// AddExp 增加经验值
func (s *UserService) AddExp(userID uint, exp int) error {
	if err := s.userRepo.UpdateExp(userID, exp); err != nil {
		return err
	}

	// 重新计算等级
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	newLevel := model.CalculateLevel(user.Exp)
	if newLevel != user.Level {
		return s.userRepo.UpdateLevel(userID, newLevel)
	}

	return nil
}

// Ban 封禁用户
func (s *UserService) Ban(userID uint) error {
	return s.userRepo.Ban(userID)
}

// Unban 解封用户
func (s *UserService) Unban(userID uint) error {
	return s.userRepo.Unban(userID)
}

// List 用户列表
func (s *UserService) List(page, size int) ([]model.User, int64, error) {
	return s.userRepo.List(page, size)
}
