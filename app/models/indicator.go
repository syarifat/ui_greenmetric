package models

import (
	"github.com/goravel/framework/database/orm"
)

type Indicator struct {
	orm.Model
	CategoryID uint     `gorm:"column:category_id"`
	Code       string   `gorm:"column:code;unique"`
	Title      string   `gorm:"column:title"`
	InputType  string   `gorm:"column:input_type"`
	MaxPoints  int      `gorm:"column:max_points"`
	Category   Category `gorm:"foreignKey:CategoryID"`
}

func (Indicator) TableName() string {
	return "indicators"
}
