package models

import (
	"time"
)

type Employeeover struct {
	Employeeoverid int       `gorm:"column:employeeoverid;primaryKey" json:"employeeoverid"`
	Companyid      int       `gorm:"column:companyid" json:"companyid"`
	Overtimeno     string    `gorm:"column:overtimeno" json:"overtimeno"`
	Overtimedate   time.Time `gorm:"column:overtimedate" json:"overtimedate"`
	Headernote     string    `gorm:"column:headernote" json:"headernote"`
	Recordstatus   int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname     string    `gorm:"column:statusname" json:"statusname"`
	Updatedate     time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company        Company   `gorm:"foreignKey:companyid;references:companyid"`
}

func (Employeeover) TableName() string {
	return "employeeover"
}
