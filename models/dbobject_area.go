package models

type DbobjectArea struct {
	ID     int    `gorm:"column:dbobjectareaid;primaryKey;autoIncrement" json:"id"`
	Name   string `gorm:"column:areaname;type:varchar(50)" json:"name"`
	X      int    `gorm:"column:sbx" json:"x"`
	Y      int    `gorm:"column:sby" json:"y"`
	Width  int    `gorm:"column:width" json:"width"`
	Height int    `gorm:"column:height" json:"height"`
	Color  string `gorm:"column:color;type:varchar(50)" json:"color"`
	Tables string `gorm:"column:tables;type:varchar(50)" json:"tables"` // JSON array of table IDs or names contained in this area
}

func (DbobjectArea) TableName() string {
	return "dbobjectarea"
}
