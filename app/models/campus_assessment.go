package models

import (
	"github.com/goravel/framework/database/orm"
)

type CampusAssessment struct {
	orm.Model
	CampusID       uint    `gorm:"column:campus_id" json:"campus_id"`
	AssessmentYear int     `gorm:"column:assessment_year" json:"assessment_year"`
	OverallScore   float64 `gorm:"column:overall_score" json:"overall_score"`
	Status         string  `gorm:"column:status" json:"status"`
	Campus         Campus  `gorm:"foreignKey:CampusID" json:"campus,omitempty"`
}

func (CampusAssessment) TableName() string {
	return "campus_assessments"
}
