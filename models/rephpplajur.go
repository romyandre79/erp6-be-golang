package models

type Rephpplajur struct {
	Companyid   int     `gorm:"column:companyid" json:"companyid"`
	Bulan       int     `gorm:"column:bulan" json:"bulan"`
	Tahun       int     `gorm:"column:tahun" json:"tahun"`
	Nourut      int     `gorm:"column:nourut" json:"nourut"`
	Keterangan  string  `gorm:"column:keterangan" json:"keterangan"`
	Plantcode   string  `gorm:"column:plantcode" json:"plantcode"`
	Accountcode string  `gorm:"column:accountcode" json:"accountcode"`
	Plantvalue  float64 `gorm:"column:plantvalue" json:"plantvalue"`
	Iscount     int8    `gorm:"column:iscount" json:"iscount"`
	Isbold      int8    `gorm:"column:isbold" json:"isbold"`
}

func (Rephpplajur) TableName() string {
	return "rephpplajur"
}
