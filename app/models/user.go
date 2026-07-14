package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	CampusID uint   `gorm:"column:campus_id" json:"campus_id"`
	Name     string `gorm:"column:name" json:"name"`
	Email    string `gorm:"column:email;unique" json:"email"`
	Password string `gorm:"column:password" json:"-"`
	Role     string `gorm:"column:role" json:"role"`
	Campus   Campus `gorm:"foreignKey:CampusID" json:"campus,omitempty"`
}

func (User) TableName() string {
	return "users"
}
