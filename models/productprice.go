package models

import (
	"time"
)

type Productprice struct {
	Productpriceid  int       `gorm:"column:productpriceid;primaryKey" json:"productpriceid"`
	Plantid         int       `gorm:"column:plantid" json:"plantid"`
	Productid       int       `gorm:"column:productid" json:"productid"`
	Buydate         time.Time `gorm:"column:buydate" json:"buydate"`
	Averageprice    float64   `gorm:"column:averageprice" json:"averageprice"`
	Oriprice        float64   `gorm:"column:oriprice" json:"oriprice"`
	Oricurrencyid   int       `gorm:"column:oricurrencyid" json:"oricurrencyid"`
	Oricurrencyrate float64   `gorm:"column:oricurrencyrate" json:"oricurrencyrate"`
	Currencyid      int       `gorm:"column:currencyid" json:"currencyid"`
	Currencyrate    float64   `gorm:"column:currencyrate" json:"currencyrate"`
	Parentrefno     string    `gorm:"column:parentrefno" json:"parentrefno"`
	Referenceno     string    `gorm:"column:referenceno" json:"referenceno"`
	Qty             float64   `gorm:"column:qty" json:"qty"`
	Uomid           int       `gorm:"column:uomid" json:"uomid"`
	Jmlhargaterima  float64   `gorm:"column:jmlhargaterima" json:"jmlhargaterima"`
	Jmlhargasisa    float64   `gorm:"column:jmlhargasisa" json:"jmlhargasisa"`
	Qtysisa         float64   `gorm:"column:qtysisa" json:"qtysisa"`
	Createddate     time.Time `gorm:"column:createddate" json:"createddate"`
	Currency        Currency  `gorm:"foreignKey:currencyid;references:currencyid"`
	Oricurrency     Currency  `gorm:"foreignKey:oricurrencyid;references:currencyid"`
	Plant           Plant     `gorm:"foreignKey:plantid;references:plantid"`
	Product         Product   `gorm:"foreignKey:productid;references:productid"`
}

func (Productprice) TableName() string {
	return "productprice"
}
