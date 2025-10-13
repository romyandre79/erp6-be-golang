package models

import (
	"time"

	"gorm.io/gorm"
)

type Language struct {
	LanguageID   uint           `gorm:"column:languageid;primaryKey;autoIncrement"`
	LanguageName string         `gorm:"column:languagename;size:20;unique;not null"`
	RecordStatus int8           `gorm:"column:recordstatus;not null;default:0"`
	UpdateDate   time.Time      `gorm:"column:updatedate;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"` // opsional kalau mau soft delete
}

// TableName override biar nama tabel sesuai
func (Language) TableName() string {
	return "language"
}
