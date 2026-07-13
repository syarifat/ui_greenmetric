package models

import (
	"github.com/goravel/framework/database/orm"
)

type Category struct {
	orm.Model
	Code             string  `gorm:"column:code;unique"`
	Name             string  `gorm:"column:name"`
	MaxPoints        int     `gorm:"column:max_points"`
	WeightPercentage float64 `gorm:"column:weight_percentage"`
}
