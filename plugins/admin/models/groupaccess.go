package admin

import (
	"time"
)

type GroupAccess struct {
	GroupAccessID int       `gorm:"column:groupaccessid;primaryKey;autoIncrement"`
	GroupName     string    `gorm:"column:groupname;size:50;unique;not null"`
	Description   string    `gorm:"column:description;size:100;not null"`
	RecordStatus  int8      `gorm:"column:recordstatus;not null"`
	UpdateDate    time.Time `gorm:"column:updatedate;autoUpdateTime"`

	// Relasi ke UserGroup (many-to-many ke UserAccess)
	UserGroups []UserGroup `gorm:"foreignKey:GroupAccessID"`
}

// Nama tabel custom
func (GroupAccess) TableName() string {
	return "groupaccess"
}
