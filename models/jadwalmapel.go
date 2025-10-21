package models

type Jadwalmapel struct {
	Jadwalmapelid int         `gorm:"column:jadwalmapelid;primaryKey" json:"jadwalmapelid"`
	Companyid     int         `gorm:"column:companyid" json:"companyid"`
	Tahunajaranid int         `gorm:"column:tahunajaranid" json:"tahunajaranid"`
	Kelasid       int         `gorm:"column:kelasid" json:"kelasid"`
	Mapelid       int         `gorm:"column:mapelid" json:"mapelid"`
	Addressbookid int         `gorm:"column:addressbookid" json:"addressbookid"`
	Monday        string      `gorm:"column:monday" json:"monday"`
	Tuesday       string      `gorm:"column:tuesday" json:"tuesday"`
	Wednesday     string      `gorm:"column:wednesday" json:"wednesday"`
	Thursday      string      `gorm:"column:thursday" json:"thursday"`
	Friday        string      `gorm:"column:friday" json:"friday"`
	Saturday      string      `gorm:"column:saturday" json:"saturday"`
	Addressbook   Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Company       Company     `gorm:"foreignKey:companyid;references:companyid"`
	Kelas         Kelas       `gorm:"foreignKey:kelasid;references:kelasid"`
	Tahunajaran   Tahunajaran `gorm:"foreignKey:tahunajaranid;references:tahunajaranid"`
}

func (Jadwalmapel) TableName() string {
	return "jadwalmapel"
}
