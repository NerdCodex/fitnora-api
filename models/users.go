package models

import "time"

type Users struct {
	UserID uint64 `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`

	UserEmail    string `gorm:"column:user_email;size:255;unique;not null" json:"user_email"`
	UserFullName string `gorm:"column:user_fullname;size:100;not null" json:"user_fullname"`
	PasswordHash string `gorm:"column:password_hash;type:text;not null" json:"-"`

	Dob    time.Time `gorm:"column:user_dob;type:date;not null" json:"user_dob"`
	Gender string    `gorm:"column:gender;size:20;not null" json:"gender"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (Users) TableName() string {
	return "users"
}
