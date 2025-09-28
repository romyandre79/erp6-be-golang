package models

import "time"

type Module struct {
	ModuleID      int       `gorm:"column:moduleid;primaryKey;autoIncrement"`
	ModuleName    string    `gorm:"column:modulename"`
	Description   string    `gorm:"column:description"`
	CreatedBy     string    `gorm:"column:createdby"`
	ModuleVersion string    `gorm:"column:moduleversion"`
	InstallDate   time.Time `gorm:"column:installdate"`
	ThemeID       *int      `gorm:"column:themeid"`
	RecordStatus  int       `gorm:"column:recordstatus"`
}

// TableName override biar GORM pakai nama tabel yang benar
func (Module) TableName() string {
	return "modules"
}
