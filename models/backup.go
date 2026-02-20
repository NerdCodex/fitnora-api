package models

import "time"

type DataBackup struct {
	DataBackupId uint64 `gorm:"column:data_backup_id;primaryKey;autoIncrement" json:"data_backup_id"`

	UserID uint64 `gorm:"column:user_id;not null;index"`
	User   Users  `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	UserDBFiles string `gorm:"column:user_dbfiles;type:text;not null" json:"-"`
	UserImages  string `gorm:"column:user_images;type:text;not null" json:"-"`

	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (DataBackup) TableName() string {
	return "data_backup"
}
