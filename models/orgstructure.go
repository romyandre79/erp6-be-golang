package models

import (
	"time"
)

type Orgstructure struct {
	Orgstructureid    int            `gorm:"column:orgstructureid;primaryKey" json:"orgstructureid"`
	Companyid         int            `gorm:"column:companyid" json:"companyid"`
	Structurename     string         `gorm:"column:structurename" json:"structurename"`
	Parentid          int            `gorm:"column:parentid" json:"parentid"`
	Amountpower       int            `gorm:"column:amountpower" json:"amountpower"`
	Startsalaryamount int            `gorm:"column:startsalaryamount" json:"startsalaryamount"`
	Endsalaryamount   int            `gorm:"column:endsalaryamount" json:"endsalaryamount"`
	Plantid           int            `gorm:"column:plantid" json:"plantid"`
	Vacancyamount     int            `gorm:"column:vacancyamount" json:"vacancyamount"`
	Isopen            int8           `gorm:"column:isopen" json:"isopen"`
	Jobdesc           string         `gorm:"column:jobdesc" json:"jobdesc"`
	Educationid       int            `gorm:"column:educationid" json:"educationid"`
	Educationmajorid  int            `gorm:"column:educationmajorid" json:"educationmajorid"`
	Startdate         time.Time      `gorm:"column:startdate" json:"startdate"`
	Enddate           time.Time      `gorm:"column:enddate" json:"enddate"`
	Recordstatus      int8           `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate        time.Time      `gorm:"column:updatedate" json:"updatedate"`
	Company           Company        `gorm:"foreignKey:companyid;references:companyid"`
	Education         Education      `gorm:"foreignKey:educationid;references:educationid"`
	Educationmajor    Educationmajor `gorm:"foreignKey:educationmajorid;references:educationmajorid"`
	Plant             Plant          `gorm:"foreignKey:plantid;references:plantid"`
}

func (Orgstructure) TableName() string {
	return "orgstructure"
}
