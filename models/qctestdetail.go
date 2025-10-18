package models

import (
  "time"
)

type Qctestdetail struct {
	Qctestdetailid int `gorm:"column:qctestdetailid;primaryKey" json:"qctestdetailid"`
	Qctestid int `gorm:"column:qctestid" json:"qctestid"`
	Qcparamid int `gorm:"column:qcparamid" json:"qcparamid"`
	Methodtest string `gorm:"column:methodtest" json:"methodtest"`
	Unittest string `gorm:"column:unittest" json:"unittest"`
	Specmin string `gorm:"column:specmin" json:"specmin"`
	Rangetarget string `gorm:"column:rangetarget" json:"rangetarget"`
	Tolerancemin string `gorm:"column:tolerancemin" json:"tolerancemin"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Qcparam Qcparam `gorm:"foreignKey:qcparamid;references:qcparamid"`
}

func (Qctestdetail) TableName() string {
	return "qctestdetail"
}
