package models

import (
	"time"
)

type Ojt struct {
	Ojtid        int       `gorm:"column:ojtid;primaryKey" json:"ojtid"`
	Plantid      int       `gorm:"column:plantid" json:"plantid"`
	Ojtno        string    `gorm:"column:ojtno" json:"ojtno"`
	Employeeid   int       `gorm:"column:employeeid" json:"employeeid"`
	Firstdate    time.Time `gorm:"column:firstdate" json:"firstdate"`
	Duedate      time.Time `gorm:"column:duedate" json:"duedate"`
	Positionid   int       `gorm:"column:positionid" json:"positionid"`
	Headernote   string    `gorm:"column:headernote" json:"headernote"`
	Recordstatus int8      `gorm:"column:recordstatus" json:"recordstatus"`
	Statusname   string    `gorm:"column:statusname" json:"statusname"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
	Employee     Employee  `gorm:"foreignKey:employeeid;references:employeeid"`
	Plant        Plant     `gorm:"foreignKey:plantid;references:plantid"`
	Position     Position  `gorm:"foreignKey:positionid;references:positionid"`
}

func (Ojt) TableName() string {
	return "ojt"
}
