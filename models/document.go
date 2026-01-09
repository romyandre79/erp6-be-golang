package models

import "time"

// Document represents the database structure for uploaded documents
type Document struct {
	DocumentID    int       `gorm:"primaryKey;column:documentid" json:"documentid"`
	UserID        int       `gorm:"column:userid" json:"userid"`
	FileName      string    `gorm:"column:filename" json:"filename"`
	FilePath      string    `gorm:"column:filepath" json:"filepath"`
	FileType      string    `gorm:"column:filetype" json:"filetype"`
	FileSize      int64     `gorm:"column:filesize" json:"filesize"`
	ExtractedText string    `gorm:"column:extractedtext;type:text" json:"extractedtext"`
	Summary       string    `gorm:"column:summary;type:text" json:"summary"`
	CreatedAt     time.Time `gorm:"column:createdat;autoCreateTime" json:"createdat"`
	UpdatedAt     time.Time `gorm:"column:updatedat;autoUpdateTime" json:"updatedat"`
}

// TableName overrides the table name used by Document to `documents`
func (Document) TableName() string {
	return "documents"
}
