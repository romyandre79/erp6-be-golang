package models

type Bank struct {
	Bankid       int    `gorm:"column:bankid;primaryKey" json:"bankid"`
	Bankname     string `gorm:"column:bankname" json:"bankname"`
	Recordstatus int8   `gorm:"column:recordstatus" json:"recordstatus"`
}

func (Bank) TableName() string {
	return "bank"
}
