package models

import (
  "time"
)

type Productoutput struct {
	Productoutputid int `gorm:"column:productoutputid;primaryKey" json:"productoutputid"`
	Productoutputno *string `gorm:"column:productoutputno" json:"productoutputno"`
	Productoutputdate time.Time `gorm:"column:productoutputdate" json:"productoutputdate"`
	Productplanid int `gorm:"column:productplanid" json:"productplanid"`
	Productplanno string `gorm:"column:productplanno" json:"productplanno"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Plantcode string `gorm:"column:plantcode" json:"plantcode"`
	Soheaderid *int `gorm:"column:soheaderid" json:"soheaderid"`
	Sono *string `gorm:"column:sono" json:"sono"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Customername *string `gorm:"column:customername" json:"customername"`
	Headernote *string `gorm:"column:headernote" json:"headernote"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname string `gorm:"column:statusname" json:"statusname"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
	Productplan Productplan `gorm:"foreignKey:productplanid;references:productplanid"`
}

func (Productoutput) TableName() string {
	return "productoutput"
}
