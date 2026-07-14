package models

import (
	"github.com/goravel/framework/database/orm"
)

type AssessmentAnswer struct {
	orm.Model
	CampusAssessmentID uint             `gorm:"column:campus_assessment_id" json:"campus_assessment_id"`
	IndicatorID        uint             `gorm:"column:indicator_id" json:"indicator_id"`
	RawInputData       string           `gorm:"column:raw_input_data" json:"raw_input_data"`
	CalculatedValue    *float64         `gorm:"column:calculated_value" json:"calculated_value"`
	SelectedTierID     *uint            `gorm:"column:selected_tier_id" json:"selected_tier_id"`
	EarnedPoints       float64          `gorm:"column:earned_points" json:"earned_points"`
	CampusAssessment   CampusAssessment `gorm:"foreignKey:CampusAssessmentID" json:"campus_assessment,omitempty"`
	Indicator          Indicator        `gorm:"foreignKey:IndicatorID" json:"indicator,omitempty"`
	Evidences          []AssessmentEvidence `gorm:"foreignKey:AssessmentAnswerID" json:"evidences,omitempty"`
}

func (AssessmentAnswer) TableName() string {
	return "assessment_answers"
}
