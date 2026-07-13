package models

import (
	"github.com/goravel/framework/database/orm"
)

type Campus struct {
	orm.Model
	Code            string `gorm:"column:code;unique"`
	Name            string `gorm:"column:name"`
	InstitutionType string `gorm:"column:institution_type"`
	Climate         string `gorm:"column:climate"`
	Setting         string `gorm:"column:setting"`
}
