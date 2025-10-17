package models

import "time"

type Groupmenu struct {
	Groupmenuid   int         `gorm:"column:groupmenuid;primaryKey" json:"groupmenuid"`
	Groupaccessid int         `gorm:"column:groupaccessid" json:"groupaccessid"`
	Menuaccessid  int         `gorm:"column:menuaccessid" json:"menuaccessid"`
	Isread        int         `gorm:"column:isread" json:"isread"`
	Iswrite       int         `gorm:"column:iswrite" json:"iswrite"`
	Ispost        int         `gorm:"column:ispost" json:"ispost"`
	Isreject      int         `gorm:"column:isreject" json:"isreject"`
	Isupload      int         `gorm:"column:isupload" json:"isupload"`
	Isdownload    int         `gorm:"column:isdownload" json:"isdownload"`
	Ispurge       int         `gorm:"column:ispurge" json:"ispurge"`
	Updatedate    time.Time   `gorm:"column:updatedate" json:"updatedate"`
	Menuaccess    Menuaccess  `gorm:"foreignKey:menuaccessid;references:menuaccessid"`
	Groupaccess   Groupaccess `gorm:"foreignKey:Groupaccessid;references:Groupaccessid" json:"groupaccess"`
}

func (Groupmenu) TableName() string {
	return "groupmenu"
}
