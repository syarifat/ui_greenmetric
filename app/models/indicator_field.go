package models

import (
	"github.com/goravel/framework/database/orm"
)

type IndicatorField struct {
	orm.Model
	IndicatorID uint      `gorm:"column:indicator_id" json:"indicator_id"`
	Key         string    `gorm:"column:key" json:"key"`
	Label       string    `gorm:"column:label" json:"label"`
	Type        string    `gorm:"column:type" json:"type"`
	Options     *string   `gorm:"column:options" json:"options"`
	Required    bool      `gorm:"column:required" json:"required"`
	Indicator   Indicator `gorm:"foreignKey:IndicatorID" json:"-"`
}

func (IndicatorField) TableName() string {
	return "indicator_fields"
}
