package models

import "time"

// ModuleTables represents the moduletables table
type ModuleTables struct {
	ModuleTablesID int       `gorm:"column:moduletablesid;primaryKey;autoIncrement" json:"moduletablesid"`
	ModuleID       int       `gorm:"column:moduleid;not null" json:"moduleid"`
	NameTable      string    `gorm:"column:nametable;type:varchar(50);not null;default:''" json:"nametable"`
	CreatedAt      time.Time `gorm:"column:createdat;not null" json:"createdat"`

	// Foreign key relationship (optional)
	Module *Modules `gorm:"foreignKey:ModuleID;references:Moduleid" json:"module,omitempty"`
}

// TableName specifies the table name for GORM
func (ModuleTables) TableName() string {
	return "moduletables"
}
