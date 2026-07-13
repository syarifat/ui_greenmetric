package models

import (
	"github.com/goravel/framework/database/orm"
)

type AssessmentAnswer struct {
	orm.Model
	CampusAssessmentID uint             `gorm:"column:campus_assessment_id"`
	IndicatorID        uint             `gorm:"column:indicator_id"`
	RawInputData       string           `gorm:"column:raw_input_data"`
	CalculatedValue    *float64         `gorm:"column:calculated_value"`
	SelectedTierID     *uint            `gorm:"column:selected_tier_id"`
	EarnedPoints       float64          `gorm:"column:earned_points"`
	CampusAssessment   CampusAssessment `gorm:"foreignKey:CampusAssessmentID"`
	Indicator          Indicator        `gorm:"foreignKey:IndicatorID"`
}
