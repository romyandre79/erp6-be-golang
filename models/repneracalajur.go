package models

type Repneracalajur struct {
	Companyid     int     `gorm:"column:companyid" json:"companyid"`
	Tahun         string  `gorm:"column:tahun" json:"tahun"`
	Bulan         int8    `gorm:"column:bulan" json:"bulan"`
	Nourut        int8    `gorm:"column:nourut" json:"nourut"`
	Accountcoded  string  `gorm:"column:accountcoded" json:"accountcoded"`
	Keterangand   string  `gorm:"column:keterangand" json:"keterangand"`
	Accountvalued float64 `gorm:"column:accountvalued" json:"accountvalued"`
	Isboldd       int8    `gorm:"column:isboldd" json:"isboldd"`
	Iscountd      int8    `gorm:"column:iscountd" json:"iscountd"`
	Accountcodek  string  `gorm:"column:accountcodek" json:"accountcodek"`
	Keterangank   string  `gorm:"column:keterangank" json:"keterangank"`
	Accountvaluek float64 `gorm:"column:accountvaluek" json:"accountvaluek"`
	Isboldk       int8    `gorm:"column:isboldk" json:"isboldk"`
	Iscountk      int8    `gorm:"column:iscountk" json:"iscountk"`
}

func (Repneracalajur) TableName() string {
	return "repneracalajur"
}
