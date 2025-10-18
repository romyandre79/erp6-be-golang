package models

import (
  "time"
)

type Cheque struct {
	Chequeid int `gorm:"column:chequeid;primaryKey" json:"chequeid"`
	Plantid int `gorm:"column:plantid" json:"plantid"`
	Isin int8 `gorm:"column:isin" json:"isin"`
	Docdate time.Time `gorm:"column:docdate" json:"docdate"`
	Bilyetgirono string `gorm:"column:bilyetgirono" json:"bilyetgirono"`
	Bankname string `gorm:"column:bankname" json:"bankname"`
	Bilyetgirovalue float64 `gorm:"column:bilyetgirovalue" json:"bilyetgirovalue"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Chequedate time.Time `gorm:"column:chequedate" json:"chequedate"`
	Expiredate time.Time `gorm:"column:expiredate" json:"expiredate"`
	Addressbookid int `gorm:"column:addressbookid" json:"addressbookid"`
	Iscair *int8 `gorm:"column:iscair" json:"iscair"`
	Cairtolakdate *time.Time `gorm:"column:cairtolakdate" json:"cairtolakdate"`
	Description *string `gorm:"column:description" json:"description"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Plant Plant `gorm:"foreignKey:plantid;references:plantid"`
}

func (Cheque) TableName() string {
	return "cheque"
}
