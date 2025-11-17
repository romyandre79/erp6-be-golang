package models

import "time"

type Menuaccess struct {
	Menuaccessid   int       `gorm:"column:menuaccessid;primaryKey" json:"menuaccessid"`
	Menuname       string    `gorm:"column:menuname" json:"menuname"`
	Menucode       string    `gorm:"column:menucode" json:"menucode"`
	Description    string    `gorm:"column:description" json:"description"`
	Moduleid       int       `gorm:"column:moduleid" json:"moduleid"`
	Parentid       *int      `gorm:"column:parentid" json:"parentid"`
	Menuurl        string    `gorm:"column:menuurl" json:"menuurl"`
	Sortorder      int       `gorm:"column:sortorder" json:"sortorder"`
	Menuicon       string    `gorm:"column:menuicon" json:"menuicon"`
	Menutype       string    `gorm:"column:menutype" json:"menutype"`
	Menuformdetail string    `gorm:"column:menuformdetail" json:"menuformdetail"`
	Recordstatus   int       `gorm:"column:recordstatus" json:"recordstatus"`
	Menuform       string    `gorm:"column:menuform" json:"menuform"`
	Updatedate     time.Time `gorm:"column:updatedate" json:"updatedate"`
	Modules        Modules   `gorm:"foreignKey:moduleid;references:moduleid"`
}

func (Menuaccess) TableName() string {
	return "menuaccess"
}
