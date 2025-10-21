package models

import (
	"time"
)

type Productoutputfg struct {
	Productoutputfgid int            `gorm:"column:productoutputfgid;primaryKey" json:"productoutputfgid"`
	Productoutputid   int            `gorm:"column:productoutputid" json:"productoutputid"`
	Productplanid     int            `gorm:"column:productplanid" json:"productplanid"`
	Productplanfgid   int            `gorm:"column:productplanfgid" json:"productplanfgid"`
	Productid         int            `gorm:"column:productid" json:"productid"`
	Productname       string         `gorm:"column:productname" json:"productname"`
	Productcode       string         `gorm:"column:productcode" json:"productcode"`
	Qty               float64        `gorm:"column:qty" json:"qty"`
	Qty2              float64        `gorm:"column:qty2" json:"qty2"`
	Qty3              float64        `gorm:"column:qty3" json:"qty3"`
	Uomid             int            `gorm:"column:uomid" json:"uomid"`
	Uom2id            int            `gorm:"column:uom2id" json:"uom2id"`
	Uom3id            int            `gorm:"column:uom3id" json:"uom3id"`
	Slocid            int            `gorm:"column:slocid" json:"slocid"`
	Storagebinid      int            `gorm:"column:storagebinid" json:"storagebinid"`
	Bomid             int            `gorm:"column:bomid" json:"bomid"`
	Mesinid           int            `gorm:"column:mesinid" json:"mesinid"`
	Processprdid      int            `gorm:"column:processprdid" json:"processprdid"`
	Employeeid        int            `gorm:"column:employeeid" json:"employeeid"`
	Shift             float64        `gorm:"column:shift" json:"shift"`
	Efektivitas       float64        `gorm:"column:efektivitas" json:"efektivitas"`
	Angkatan          float64        `gorm:"column:angkatan" json:"angkatan"`
	Outputdate        time.Time      `gorm:"column:outputdate" json:"outputdate"`
	Lotno             string         `gorm:"column:lotno" json:"lotno"`
	Description       string         `gorm:"column:description" json:"description"`
	Nilaihpp          float64        `gorm:"column:nilaihpp" json:"nilaihpp"`
	Qtyres            float64        `gorm:"column:qtyres" json:"qtyres"`
	Tsqty             float64        `gorm:"column:tsqty" json:"tsqty"`
	Tsqty2            float64        `gorm:"column:tsqty2" json:"tsqty2"`
	Tsqty3            float64        `gorm:"column:tsqty3" json:"tsqty3"`
	Updatedate        time.Time      `gorm:"column:updatedate" json:"updatedate"`
	Billofmaterial    Billofmaterial `gorm:"foreignKey:bomid;references:bomid"`
	Employee          Employee       `gorm:"foreignKey:employeeid;references:employeeid"`
	Mesin             Mesin          `gorm:"foreignKey:mesinid;references:mesinid"`
	Productplan       Productplan    `gorm:"foreignKey:productplanid;references:productplanid"`
	Productplanfg     Productplanfg  `gorm:"foreignKey:productplanfgid;references:productplanfgid"`
	Product           Product        `gorm:"foreignKey:productid;references:productid"`
	Sloc              Sloc           `gorm:"foreignKey:slocid;references:slocid"`
	Uom               Unitofmeasure  `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2              Unitofmeasure  `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3              Unitofmeasure  `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
}

func (Productoutputfg) TableName() string {
	return "productoutputfg"
}
