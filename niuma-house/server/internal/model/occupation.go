package model

import "time"

// Occupation 职业分类实体
type Occupation struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (Occupation) TableName() string {
	return "occupations"
}

// 预置职业分类
var DefaultOccupations = []Occupation{
	{ID: 1, Name: "程序员"},
	{ID: 2, Name: "产品经理"},
	{ID: 3, Name: "运营"},
	{ID: 4, Name: "设计师"},
	{ID: 5, Name: "销售"},
	{ID: 6, Name: "人事HR"},
	{ID: 7, Name: "财务"},
	{ID: 8, Name: "市场"},
	{ID: 9, Name: "其他"},
}
