package models

import (
	"github.com/goravel/framework/database/orm"
)

type Campus struct {
	orm.Model
	Code            string `gorm:"column:code;unique" json:"code"`
	Name            string `gorm:"column:name" json:"name"`
	InstitutionType string `gorm:"column:institution_type" json:"institution_type"`
	Climate         string `gorm:"column:climate" json:"climate"`
	Setting         string `gorm:"column:setting" json:"setting"`
}

func (Campus) TableName() string {
	return "campuses"
}
