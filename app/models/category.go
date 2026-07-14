package models

import (
	"github.com/goravel/framework/database/orm"
)

type Category struct {
	orm.Model
	Code             string  `gorm:"column:code;unique" json:"code"`
	Name             string  `gorm:"column:name" json:"name"`
	MaxPoints        int     `gorm:"column:max_points" json:"max_points"`
	WeightPercentage float64 `gorm:"column:weight_percentage" json:"weight_percentage"`
}

func (Category) TableName() string {
	return "categories"
}
