package models

import (
	"github.com/goravel/framework/database/orm"
)

type IndicatorScoringTier struct {
	orm.Model
	IndicatorID     uint      `gorm:"column:indicator_id"`
	OptionLabel     string    `gorm:"column:option_label"`
	MinValue        *float64  `gorm:"column:min_value"`
	MaxValue        *float64  `gorm:"column:max_value"`
	Operator        string    `gorm:"column:operator"`
	PointMultiplier float64   `gorm:"column:point_multiplier"`
	Indicator       Indicator `gorm:"foreignKey:IndicatorID"`
}

func (IndicatorScoringTier) TableName() string {
	return "indicator_scoring_tiers"
}
