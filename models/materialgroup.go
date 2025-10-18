package models

import (
  "time"
)

type Materialgroup struct {
	Materialgroupid int `gorm:"column:materialgroupid;primaryKey" json:"materialgroupid"`
	Materialgroupcode string `gorm:"column:materialgroupcode" json:"materialgroupcode"`
	Description string `gorm:"column:description" json:"description"`
	Parentmatgroupid *int `gorm:"column:parentmatgroupid" json:"parentmatgroupid"`
	Isfg int8 `gorm:"column:isfg" json:"isfg"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Materialgroup) TableName() string {
	return "materialgroup"
}
