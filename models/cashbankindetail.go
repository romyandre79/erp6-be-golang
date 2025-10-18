package models

import (
  "time"
)

type Cashbankindetail struct {
	Cashbankindetailid int `gorm:"column:cashbankindetailid;primaryKey" json:"cashbankindetailid"`
	Cashbankinid int `gorm:"column:cashbankinid" json:"cashbankinid"`
	Invoicearid int `gorm:"column:invoicearid" json:"invoicearid"`
	Amount float64 `gorm:"column:amount" json:"amount"`
	Currencyid int `gorm:"column:currencyid" json:"currencyid"`
	Ratevalue float64 `gorm:"column:ratevalue" json:"ratevalue"`
	Detailnote string `gorm:"column:detailnote" json:"detailnote"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Currency Currency `gorm:"foreignKey:currencyid;references:currencyid"`
	Invoicear Invoicear `gorm:"foreignKey:invoicearid;references:invoicearid"`
}

func (Cashbankindetail) TableName() string {
	return "cashbankindetail"
}
