package models

import (
  "time"
)

type Productstockdet struct {
	Productstockdetid int `gorm:"column:productstockdetid;primaryKey" json:"productstockdetid"`
	Productid int `gorm:"column:productid" json:"productid"`
	Addressbookid *int `gorm:"column:addressbookid" json:"addressbookid"`
	Productcode *string `gorm:"column:productcode" json:"productcode"`
	Productname *string `gorm:"column:productname" json:"productname"`
	Qty float64 `gorm:"column:qty" json:"qty"`
	Qty2 float64 `gorm:"column:qty2" json:"qty2"`
	Qty3 *float64 `gorm:"column:qty3" json:"qty3"`
	Uomid int `gorm:"column:uomid" json:"uomid"`
	Uomcode string `gorm:"column:uomcode" json:"uomcode"`
	Uom2id int `gorm:"column:uom2id" json:"uom2id"`
	Uom2code string `gorm:"column:uom2code" json:"uom2code"`
	Uom3id *int `gorm:"column:uom3id" json:"uom3id"`
	Uom3code *string `gorm:"column:uom3code" json:"uom3code"`
	Slocid int `gorm:"column:slocid" json:"slocid"`
	Sloccode string `gorm:"column:sloccode" json:"sloccode"`
	Referenceno *string `gorm:"column:referenceno" json:"referenceno"`
	Productstockid *int `gorm:"column:productstockid" json:"productstockid"`
	Storagebinid int `gorm:"column:storagebinid" json:"storagebinid"`
	Storagedesc string `gorm:"column:storagedesc" json:"storagedesc"`
	Transdate *time.Time `gorm:"column:transdate" json:"transdate"`
	Transstockdetid *int `gorm:"column:transstockdetid" json:"transstockdetid"`
	Formrequestrawid *int `gorm:"column:formrequestrawid" json:"formrequestrawid"`
	Formrequestresultid *int `gorm:"column:formrequestresultid" json:"formrequestresultid"`
	Productoutputfgid *int `gorm:"column:productoutputfgid" json:"productoutputfgid"`
	Productoutputdetailid *int `gorm:"column:productoutputdetailid" json:"productoutputdetailid"`
	Transstockid *int `gorm:"column:transstockid" json:"transstockid"`
	Expiredate *time.Time `gorm:"column:expiredate" json:"expiredate"`
	Buydate *time.Time `gorm:"column:buydate" json:"buydate"`
	Buyprice *float64 `gorm:"column:buyprice" json:"buyprice"`
	Currencyid *int `gorm:"column:currencyid" json:"currencyid"`
	Currencyrate *float64 `gorm:"column:currencyrate" json:"currencyrate"`
	Materialstatusid *int `gorm:"column:materialstatusid" json:"materialstatusid"`
	Materialstatusname *string `gorm:"column:materialstatusname" json:"materialstatusname"`
	Ownershipid *int `gorm:"column:ownershipid" json:"ownershipid"`
	Ownershipname *string `gorm:"column:ownershipname" json:"ownershipname"`
	Lotno *string `gorm:"column:lotno" json:"lotno"`
	Grdetailid *int `gorm:"column:grdetailid" json:"grdetailid"`
	Qtyalokasi *float64 `gorm:"column:qtyalokasi" json:"qtyalokasi"`
	Alokasifrrawid *int `gorm:"column:alokasifrrawid" json:"alokasifrrawid"`
	Alokasiqty *float64 `gorm:"column:alokasiqty" json:"alokasiqty"`
	Averageprice *float64 `gorm:"column:averageprice" json:"averageprice"`
	Vrqty *float64 `gorm:"column:vrqty" json:"vrqty"`
	Vrqty2 *float64 `gorm:"column:vrqty2" json:"vrqty2"`
	Vrqty3 *float64 `gorm:"column:vrqty3" json:"vrqty3"`
	Product Product `gorm:"foreignKey:productid;references:productid"`
	Storagebin Storagebin `gorm:"foreignKey:storagebinid;references:storagebinid"`
	Sloc Sloc `gorm:"foreignKey:slocid;references:slocid"`
	Unitofmeasure Unitofmeasure `gorm:"foreignKey:uomid;references:unitofmeasureid"`
}

func (Productstockdet) TableName() string {
	return "productstockdet"
}
