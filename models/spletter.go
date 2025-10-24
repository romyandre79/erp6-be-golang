package models

import (
	"time"
)

type Spletter struct {
	Spletterid     int          `gorm:"column:spletterid;primaryKey" json:"spletterid"`
	Plantid        int          `gorm:"column:plantid" json:"plantid"`
	Spletterno     string       `gorm:"column:spletterno" json:"spletterno"`
	Employeeid     int          `gorm:"column:employeeid" json:"employeeid"`
	Splettertypeid int          `gorm:"column:splettertypeid" json:"splettertypeid"`
	Spletterdate   time.Time    `gorm:"column:spletterdate" json:"spletterdate"`
	Description    string       `gorm:"column:description" json:"description"`
	Recordstatus   int8         `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname     string       `gorm:"column:statusname" json:"statusname"`
	Updatedate     time.Time    `gorm:"column:updatedate" json:"updatedate"`
	Employee       Employee     `gorm:"foreignKey:employeeid;references:employeeid"`
	Plant          Plant        `gorm:"foreignKey:plantid;references:plantid"`
	Splettertype   Splettertype `gorm:"foreignKey:splettertypeid;references:splettertypeid"`
}

func (Spletter) TableName() string {
	return "spletter"
}
