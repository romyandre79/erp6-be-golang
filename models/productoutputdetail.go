package models

import (
	"time"
)

type Productoutputdetail struct {
	Productoutputdetailid int            `gorm:"column:productoutputdetailid;primaryKey" json:"productoutputdetailid"`
	Productoutputid       int            `gorm:"column:productoutputid" json:"productoutputid"`
	Productoutputfgid     int            `gorm:"column:productoutputfgid" json:"productoutputfgid"`
	Productid             int            `gorm:"column:productid" json:"productid"`
	Productcode           string         `gorm:"column:productcode" json:"productcode"`
	Productname           string         `gorm:"column:productname" json:"productname"`
	Qty                   float64        `gorm:"column:qty" json:"qty"`
	Qty2                  float64        `gorm:"column:qty2" json:"qty2"`
	Qty3                  float64        `gorm:"column:qty3" json:"qty3"`
	Uomid                 int            `gorm:"column:uomid" json:"uomid"`
	Uom2id                int            `gorm:"column:uom2id" json:"uom2id"`
	Uom3id                int            `gorm:"column:uom3id" json:"uom3id"`
	Slocfromid            int            `gorm:"column:slocfromid" json:"slocfromid"`
	Sloctoid              int            `gorm:"column:sloctoid" json:"sloctoid"`
	Storagebintoid        int            `gorm:"column:storagebintoid" json:"storagebintoid"`
	Productplandetailid   int            `gorm:"column:productplandetailid" json:"productplandetailid"`
	Productplanfgid       int            `gorm:"column:productplanfgid" json:"productplanfgid"`
	Outputdate            time.Time      `gorm:"column:outputdate" json:"outputdate"`
	Bomid                 int            `gorm:"column:bomid" json:"bomid"`
	Description           string         `gorm:"column:description" json:"description"`
	Nilaihpp              float64        `gorm:"column:nilaihpp" json:"nilaihpp"`
	Updatedate            time.Time      `gorm:"column:updatedate" json:"updatedate"`
	Billofmaterial        Billofmaterial `gorm:"foreignKey:bomid;references:bomid"`
	Product               Product        `gorm:"foreignKey:productid;references:productid"`
	Slocfrom              Sloc           `gorm:"foreignKey:slocfromid;references:slocid"`
	Slocto                Sloc           `gorm:"foreignKey:sloctoid;references:slocid"`
	Unitofmeasure         Unitofmeasure  `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Productoutputdetail) TableName() string {
	return "productoutputdetail"
}
