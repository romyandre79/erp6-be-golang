package models

type Repbiayaumumlajur struct {
	Companyid   int    `gorm:"column:companyid" json:"companyid"`
	Bulan       int    `gorm:"column:bulan" json:"bulan"`
	Tahun       int    `gorm:"column:tahun" json:"tahun"`
	Nourut      string `gorm:"column:nourut" json:"nourut"`
	Keterangan  string `gorm:"column:keterangan" json:"keterangan"`
	Accountcode string `gorm:"column:accountcode" json:"accountcode"`
	Plantcode   string `gorm:"column:plantcode" json:"plantcode"`
	Plantvalue  string `gorm:"column:plantvalue" json:"plantvalue"`
	Isbold      int8   `gorm:"column:isbold" json:"isbold"`
	Iscount     int8   `gorm:"column:iscount" json:"iscount"`
}

func (Repbiayaumumlajur) TableName() string {
	return "repbiayaumumlajur"
}
