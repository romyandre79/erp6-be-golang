package models

type Walimurid struct {
	Walimuridid   int         `gorm:"column:walimuridid;primaryKey" json:"walimuridid"`
	Addressbookid int         `gorm:"column:addressbookid" json:"addressbookid"`
	Siswaid       int         `gorm:"column:siswaid" json:"siswaid"`
	Addressbook   Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Siswa         Addressbook `gorm:"foreignKey:siswaid;references:addressbookid"`
}

func (Walimurid) TableName() string {
	return "walimurid"
}
