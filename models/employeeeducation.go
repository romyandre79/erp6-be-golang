package models

import (
  "time"
)

type Employeeeducation struct {
	Employeeeducationid int `gorm:"column:employeeeducationid;primaryKey" json:"employeeeducationid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Educationid int `gorm:"column:educationid" json:"educationid"`
	Schoolname string `gorm:"column:schoolname" json:"schoolname"`
	Cityid int `gorm:"column:cityid" json:"cityid"`
	Graduate *string `gorm:"column:graduate" json:"graduate"`
	Isdiploma int8 `gorm:"column:isdiploma" json:"isdiploma"`
	Schooldegree string `gorm:"column:schooldegree" json:"schooldegree"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	City City `gorm:"foreignKey:cityid;references:cityid"`
	Education Education `gorm:"foreignKey:educationid;references:educationid"`
}

func (Employeeeducation) TableName() string {
	return "employeeeducation"
}
