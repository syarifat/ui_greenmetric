package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	CampusID uint   `gorm:"column:campus_id"`
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email;unique"`
	Password string `gorm:"column:password"`
	Role     string `gorm:"column:role"`
	Campus   Campus `gorm:"foreignKey:CampusID"`
}

func (User) TableName() string {
	return "users"
}
