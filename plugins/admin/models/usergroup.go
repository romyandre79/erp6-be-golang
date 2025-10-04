package admin

import (
	"time"
)

type UserGroup struct {
	UserGroupID   int       `gorm:"column:usergroupid;primaryKey;autoIncrement"`
	UserAccessID  int       `gorm:"column:useraccessid;not null"`
	GroupAccessID int       `gorm:"column:groupaccessid;not null"`
	UpdateDate    time.Time `gorm:"column:updatedate;autoUpdateTime"`

	// Relasi
	GroupAccess GroupAccess `gorm:"foreignKey:GroupAccessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserAccess  UserAccess  `gorm:"foreignKey:UserAccessID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// Nama tabel custom
func (UserGroup) TableName() string {
	return "usergroup"
}
