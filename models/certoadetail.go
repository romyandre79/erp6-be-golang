package models

import (
  "time"
)

type Certoadetail struct {
	Certoadetailid int `gorm:"column:certoadetailid;primaryKey" json:"certoadetailid"`
	Certoaid int `gorm:"column:certoaid" json:"certoaid"`
	Qcparamid int `gorm:"column:qcparamid" json:"qcparamid"`
	Methodtest string `gorm:"column:methodtest" json:"methodtest"`
	Unittest string `gorm:"column:unittest" json:"unittest"`
	Specmin string `gorm:"column:specmin" json:"specmin"`
	Rangetarget string `gorm:"column:rangetarget" json:"rangetarget"`
	Tolerancemin string `gorm:"column:tolerancemin" json:"tolerancemin"`
	Resulttest string `gorm:"column:resulttest" json:"resulttest"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Qcparam Qcparam `gorm:"foreignKey:qcparamid;references:qcparamid"`
}

func (Certoadetail) TableName() string {
	return "certoadetail"
}
