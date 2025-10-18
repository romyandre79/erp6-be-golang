package models

import (
  "time"
)

type Employeefamily struct {
	Employeefamilyid int `gorm:"column:employeefamilyid;primaryKey" json:"employeefamilyid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Familyrelationid int `gorm:"column:familyrelationid" json:"familyrelationid"`
	Familyname string `gorm:"column:familyname" json:"familyname"`
	Sexid int `gorm:"column:sexid" json:"sexid"`
	Cityid int `gorm:"column:cityid" json:"cityid"`
	Birthdate time.Time `gorm:"column:birthdate" json:"birthdate"`
	Educationid int `gorm:"column:educationid" json:"educationid"`
	Occupationid int `gorm:"column:occupationid" json:"occupationid"`
	Bpjskesno string `gorm:"column:bpjskesno" json:"bpjskesno"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	City City `gorm:"foreignKey:cityid;references:cityid"`
	Education Education `gorm:"foreignKey:educationid;references:educationid"`
	Familyrelation Familyrelation `gorm:"foreignKey:familyrelationid;references:familyrelationid"`
	Sex Sex `gorm:"foreignKey:sexid;references:sexid"`
}

func (Employeefamily) TableName() string {
	return "employeefamily"
}
