package models

import (
  "time"
)

type Bookrent struct {
	Bookrentid int `gorm:"column:bookrentid;primaryKey" json:"bookrentid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Bookid int `gorm:"column:bookid" json:"bookid"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Backdate *time.Time `gorm:"column:backdate" json:"backdate"`
	Price float64 `gorm:"column:price" json:"price"`
	Pricelate *float64 `gorm:"column:pricelate" json:"pricelate"`
	Isreturn *int8 `gorm:"column:isreturn" json:"isreturn"`
	Duedate time.Time `gorm:"column:duedate" json:"duedate"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate time.Time `gorm:"column:createddate" json:"createddate"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Addressbook Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Book Book `gorm:"foreignKey:bookid;references:bookid"`
	Company Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Bookrent) TableName() string {
	return "bookrent"
}
