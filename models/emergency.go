package models

import (
  "time"
)

type Emergency struct {
	Emergencyid int `gorm:"column:emergencyid;primaryKey" json:"emergencyid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Emergencydate time.Time `gorm:"column:emergencydate" json:"emergencydate"`
	Emergencyno string `gorm:"column:emergencyno" json:"emergencyno"`
	Orgstructureid int `gorm:"column:orgstructureid" json:"orgstructureid"`
	Headernote *string `gorm:"column:headernote" json:"headernote"`
	Recordstatus int `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Orgstructure Orgstructure `gorm:"foreignKey:orgstructureid;references:orgstructureid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Emergency) TableName() string {
	return "emergency"
}
