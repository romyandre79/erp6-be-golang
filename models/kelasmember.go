package models

import (
  "time"
)

type Kelasmember struct {
	Kelasmemberid int `gorm:"column:kelasmemberid;primaryKey" json:"kelasmemberid"`
	Kelasid int `gorm:"column:kelasid" json:"kelasid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Tahunajaranid int `gorm:"column:tahunajaranid" json:"tahunajaranid"`
	Statusmemberid int `gorm:"column:statusmemberid" json:"statusmemberid"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Kelas Kelas `gorm:"foreignKey:kelasid;references:kelasid"`
	Statusmember Statusmember `gorm:"foreignKey:statusmemberid;references:statusmemberid"`
	Tahunajaran Tahunajaran `gorm:"foreignKey:tahunajaranid;references:tahunajaranid"`
}

func (Kelasmember) TableName() string {
	return "kelasmember"
}
