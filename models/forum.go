package models

import (
	"time"
)

type Forum struct {
	Forumid       int         `gorm:"column:forumid;primaryKey" json:"forumid"`
	Forumtitle    string      `gorm:"column:forumtitle" json:"forumtitle"`
	Forumcontent  string      `gorm:"column:forumcontent" json:"forumcontent"`
	Addressbookid int         `gorm:"column:addressbookid" json:"addressbookid"`
	Replyid       int         `gorm:"column:replyid" json:"replyid"`
	Postdate      time.Time   `gorm:"column:postdate" json:"postdate"`
	Recordstatus  int8        `gorm:"column:recordstatus" json:"recordstatus"`
	Addressbook   Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Replyforum    *Forum      `gorm:"foreignKey:replyid;references:forumid"`
}

func (Forum) TableName() string {
	return "forum"
}
