package models

import (
	"time"
)

type Transstockdet struct {
	Transstockdetid     int             `gorm:"column:transstockdetid;primaryKey" json:"transstockdetid"`
	Transstockid        int             `gorm:"column:transstockid" json:"transstockid"`
	Formrequestrawid    int             `gorm:"column:formrequestrawid" json:"formrequestrawid"`
	Formrequestresultid int             `gorm:"column:formrequestresultid" json:"formrequestresultid"`
	Productplandetailid int             `gorm:"column:productplandetailid" json:"productplandetailid"`
	Productplanfgid     int             `gorm:"column:productplanfgid" json:"productplanfgid"`
	Productid           int             `gorm:"column:productid" json:"productid"`
	Productcode         string          `gorm:"column:productcode" json:"productcode"`
	Productname         string          `gorm:"column:productname" json:"productname"`
	Qty                 float64         `gorm:"column:qty" json:"qty"`
	Qtyreq              float64         `gorm:"column:qtyreq" json:"qtyreq"`
	Qty2                float64         `gorm:"column:qty2" json:"qty2"`
	Qty3                float64         `gorm:"column:qty3" json:"qty3"`
	Uomid               int             `gorm:"column:uomid" json:"uomid"`
	Uom2id              int             `gorm:"column:uom2id" json:"uom2id"`
	Uom3id              int             `gorm:"column:uom3id" json:"uom3id"`
	Mesinid             int             `gorm:"column:mesinid" json:"mesinid"`
	Storagebinfromid    int             `gorm:"column:storagebinfromid" json:"storagebinfromid"`
	Storagebintoid      int             `gorm:"column:storagebintoid" json:"storagebintoid"`
	Productoutputfgid   int             `gorm:"column:productoutputfgid" json:"productoutputfgid"`
	Lotno               string          `gorm:"column:lotno" json:"lotno"`
	Itemnote            string          `gorm:"column:itemnote" json:"itemnote"`
	Updatedate          time.Time       `gorm:"column:updatedate" json:"updatedate"`
	Mesin               Mesin           `gorm:"foreignKey:mesinid;references:mesinid"`
	Product             Product         `gorm:"foreignKey:productid;references:productid"`
	Productoutputfg     Productoutputfg `gorm:"foreignKey:productoutputfgid;references:productoutputfgid"`
	Storagebinfrom      Storagebin      `gorm:"foreignKey:storagebinfromid;references:storagebinid"`
	Storagebinto        Storagebin      `gorm:"foreignKey:storagebintoid;references:storagebinid"`
	Uom                 Unitofmeasure   `gorm:"foreignKey:uomid;references:unitofmeasureid"`
	Uom2                Unitofmeasure   `gorm:"foreignKey:uom2id;references:unitofmeasureid"`
	Uom3                Unitofmeasure   `gorm:"foreignKey:uom3id;references:unitofmeasureid"`
}

func (Transstockdet) TableName() string {
	return "transstockdet"
}
