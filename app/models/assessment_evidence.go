package models

import (
	"github.com/goravel/framework/database/orm"
)

type AssessmentEvidence struct {
	orm.Model
	AssessmentAnswerID uint             `gorm:"column:assessment_answer_id" json:"assessment_answer_id"`
	DocumentName       string           `gorm:"column:document_name" json:"document_name"`
	Description        *string          `gorm:"column:description" json:"description"`
	FileUrl            string           `gorm:"column:file_url" json:"file_url"`
	AssessmentAnswer   AssessmentAnswer `gorm:"foreignKey:AssessmentAnswerID" json:"assessment_answer,omitempty"`
}

func (AssessmentEvidence) TableName() string {
	return "assessment_evidences"
}
