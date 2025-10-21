package models

import (
	"time"
)

type Addressbook struct {
	Addressbookid    int           `gorm:"column:addressbookid;primaryKey" json:"addressbookid"`
	Fullname         string        `gorm:"column:fullname" json:"fullname"`
	Addressbookno    string        `gorm:"column:addressbookno" json:"addressbookno"`
	Iscustomer       int8          `gorm:"column:iscustomer" json:"iscustomer"`
	Isapplicant      int8          `gorm:"column:isapplicant" json:"isapplicant"`
	Isemployee       int8          `gorm:"column:isemployee" json:"isemployee"`
	Isvendor         int8          `gorm:"column:isvendor" json:"isvendor"`
	Ishospital       int8          `gorm:"column:ishospital" json:"ishospital"`
	Isexpedisi       int8          `gorm:"column:isexpedisi" json:"isexpedisi"`
	Ismember         int8          `gorm:"column:ismember" json:"ismember"`
	Isguru           int8          `gorm:"column:isguru" json:"isguru"`
	Iskasek          int8          `gorm:"column:iskasek" json:"iskasek"`
	Currentlimit     float64       `gorm:"column:currentlimit" json:"currentlimit"`
	Currentdebt      float64       `gorm:"column:currentdebt" json:"currentdebt"`
	Taxno            string        `gorm:"column:taxno" json:"taxno"`
	Jenisanggotaid   int           `gorm:"column:jenisanggotaid" json:"jenisanggotaid"`
	Creditlimit      float64       `gorm:"column:creditlimit" json:"creditlimit"`
	Isstrictlimit    int8          `gorm:"column:isstrictlimit" json:"isstrictlimit"`
	Salesareaid      int           `gorm:"column:salesareaid" json:"salesareaid"`
	Pricecategoryid  int           `gorm:"column:pricecategoryid" json:"pricecategoryid"`
	Overdue          int           `gorm:"column:overdue" json:"overdue"`
	Invoicedate      time.Time     `gorm:"column:invoicedate" json:"invoicedate"`
	Top              string        `gorm:"column:top" json:"top"`
	Hutang           float64       `gorm:"column:hutang" json:"hutang"`
	Pendinganso      float64       `gorm:"column:pendinganso" json:"pendinganso"`
	Logo             string        `gorm:"column:logo" json:"logo"`
	Groupcustomerid  int           `gorm:"column:groupcustomerid" json:"groupcustomerid"`
	Paymentmethodid  int           `gorm:"column:paymentmethodid" json:"paymentmethodid"`
	Taxid            int           `gorm:"column:taxid" json:"taxid"`
	Companytypeid    int           `gorm:"column:companytypeid" json:"companytypeid"`
	Pin              string        `gorm:"column:pin" json:"pin"`
	Useraccessid     int           `gorm:"column:useraccessid" json:"useraccessid"`
	Website          string        `gorm:"column:website" json:"website"`
	Katamutiara      string        `gorm:"column:katamutiara" json:"katamutiara"`
	Positionid       int           `gorm:"column:positionid" json:"positionid"`
	Employeetypeid   int           `gorm:"column:employeetypeid" json:"employeetypeid"`
	Kutypeid         int           `gorm:"column:kutypeid" json:"kutypeid"`
	Sexid            int           `gorm:"column:sexid" json:"sexid"`
	Birthcityid      int           `gorm:"column:birthcityid" json:"birthcityid"`
	Birthdate        time.Time     `gorm:"column:birthdate" json:"birthdate"`
	Religionid       int           `gorm:"column:religionid" json:"religionid"`
	Maritalstatusid  int           `gorm:"column:maritalstatusid" json:"maritalstatusid"`
	Referenceby      string        `gorm:"column:referenceby" json:"referenceby"`
	Joindate         time.Time     `gorm:"column:joindate" json:"joindate"`
	Employeestatusid int           `gorm:"column:employeestatusid" json:"employeestatusid"`
	Istrial          int8          `gorm:"column:istrial" json:"istrial"`
	Resigndate       time.Time     `gorm:"column:resigndate" json:"resigndate"`
	Levelorgid       int           `gorm:"column:levelorgid" json:"levelorgid"`
	Recordstatus     int8          `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname       string        `gorm:"column:statusname" json:"statusname"`
	Updatedate       time.Time     `gorm:"column:updatedate" json:"updatedate"`
	City             City          `gorm:"foreignKey:birthcityid;references:cityid"`
	Companytype      Companytype   `gorm:"foreignKey:companytypeid;references:companytypeid"`
	Employeetype     Employeetype  `gorm:"foreignKey:employeetypeid;references:employeetypeid"`
	Groupcustomer    Groupcustomer `gorm:"foreignKey:groupcustomerid;references:groupcustomerid"`
	Jenisanggota     Jenisanggota  `gorm:"foreignKey:jenisanggotaid;references:jenisanggotaid"`
	Kutype           Kutype        `gorm:"foreignKey:kutypeid;references:kutypeid"`
	Paymentmethod    Paymentmethod `gorm:"foreignKey:paymentmethodid;references:paymentmethodid"`
	Position         Position      `gorm:"foreignKey:positionid;references:positionid"`
	Pricecategory    Pricecategory `gorm:"foreignKey:pricecategoryid;references:pricecategoryid"`
	Religion         Religion      `gorm:"foreignKey:religionid;references:religionid"`
	Salesarea        Salesarea     `gorm:"foreignKey:salesareaid;references:salesareaid"`
	Sex              Sex           `gorm:"foreignKey:sexid;references:sexid"`
	Tax              Tax           `gorm:"foreignKey:taxid;references:taxid"`
	Useraccess       Useraccess    `gorm:"foreignKey:useraccessid;references:useraccessid"`
}

func (Addressbook) TableName() string {
	return "addressbook"
}
