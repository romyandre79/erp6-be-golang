package admin

import (
	"time"
)

type UserGroup struct {
	UserGroupID   uint      `gorm:"primaryKey;column:usergroupid"`
	UserAccessID  uint      `gorm:"column:useraccessid"`
	GroupAccessID uint      `gorm:"column:groupaccessid"`
	UpdateDate    time.Time `gorm:"column:updatedate"`

	// Relasi eksplisit, biar GORM tahu arah referensinya
	GroupAccess GroupAccess `gorm:"foreignKey:GroupAccessID;references:GroupAccessID"`
	UserAccess  UserAccess  `gorm:"foreignKey:UserAccessID;references:UserAccessID"`
}

// Nama tabel custom
func (UserGroup) TableName() string {
	return "usergroup"
}
