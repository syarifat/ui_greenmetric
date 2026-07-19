package models

import (
	"github.com/goravel/framework/database/orm"
)

type Indicator struct {
	orm.Model
	CategoryID uint     `gorm:"column:category_id" json:"category_id"`
	Code       string   `gorm:"column:code;unique" json:"code"`
	Title      string   `gorm:"column:title" json:"title"`
	InputType  string   `gorm:"column:input_type" json:"input_type"`
	MaxPoints  int      `gorm:"column:max_points" json:"max_points"`
	Category   Category         `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Fields     []IndicatorField `gorm:"foreignKey:IndicatorID" json:"fields,omitempty"`
}

func (Indicator) TableName() string {
	return "indicators"
}
