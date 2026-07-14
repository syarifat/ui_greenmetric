package models

import (
	"github.com/goravel/framework/database/orm"
)

type IndicatorScoringTier struct {
	orm.Model
	IndicatorID     uint      `gorm:"column:indicator_id" json:"indicator_id"`
	OptionLabel     string    `gorm:"column:option_label" json:"option_label"`
	MinValue        *float64  `gorm:"column:min_value" json:"min_value"`
	MaxValue        *float64  `gorm:"column:max_value" json:"max_value"`
	Operator        string    `gorm:"column:operator" json:"operator"`
	PointMultiplier float64   `gorm:"column:point_multiplier" json:"point_multiplier"`
	Indicator       Indicator `gorm:"foreignKey:IndicatorID" json:"indicator,omitempty"`
}

func (IndicatorScoringTier) TableName() string {
	return "indicator_scoring_tiers"
}
