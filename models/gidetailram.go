package models

type Gidetailram struct {
	Gidetailid        int     `gorm:"column:gidetailid" json:"gidetailid"`
	Giheaderid        int     `gorm:"column:giheaderid" json:"giheaderid"`
	Productstockdetid int     `gorm:"column:productstockdetid" json:"productstockdetid"`
	Productid         int     `gorm:"column:productid" json:"productid"`
	Qty               float64 `gorm:"column:qty" json:"qty"`
	Uomid             int     `gorm:"column:uomid" json:"uomid"`
	Qty2              float64 `gorm:"column:qty2" json:"qty2"`
	Uom2id            int     `gorm:"column:uom2id" json:"uom2id"`
	Qty3              float64 `gorm:"column:qty3" json:"qty3"`
	Uom3id            int     `gorm:"column:uom3id" json:"uom3id"`
	Qty4              float64 `gorm:"column:qty4" json:"qty4"`
	Uom4id            float64 `gorm:"column:uom4id" json:"uom4id"`
	Slocid            int     `gorm:"column:slocid" json:"slocid"`
	Itemnote          string  `gorm:"column:itemnote" json:"itemnote"`
	Serialno          string  `gorm:"column:serialno" json:"serialno"`
	Sodetailid        int     `gorm:"column:sodetailid" json:"sodetailid"`
	Storagebinid      int     `gorm:"column:storagebinid" json:"storagebinid"`
}

func (Gidetailram) TableName() string {
	return "gidetailram"
}
