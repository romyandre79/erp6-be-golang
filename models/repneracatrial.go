package models

type Repneracatrial struct {
	Companyid    int     `gorm:"column:companyid" json:"companyid"`
	Tahun        int     `gorm:"column:tahun" json:"tahun"`
	Bulan        int     `gorm:"column:bulan" json:"bulan"`
	Nourut       int     `gorm:"column:nourut" json:"nourut"`
	Accountcode  string  `gorm:"column:accountcode" json:"accountcode"`
	Keterangan   string  `gorm:"column:keterangan" json:"keterangan"`
	Accountvalue float64 `gorm:"column:accountvalue" json:"accountvalue"`
}

func (Repneracatrial) TableName() string {
	return "repneracatrial"
}
