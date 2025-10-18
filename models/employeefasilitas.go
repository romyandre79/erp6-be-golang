package models

import (
  "time"
)

type Employeefasilitas struct {
	Employeefasilitasid int `gorm:"column:employeefasilitasid;primaryKey" json:"employeefasilitasid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Fixassetid int `gorm:"column:fixassetid" json:"fixassetid"`
	Description string `gorm:"column:description" json:"description"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Fixasset Fixasset `gorm:"foreignKey:fixassetid;references:fixassetid"`
}

func (Employeefasilitas) TableName() string {
	return "employeefasilitas"
}
