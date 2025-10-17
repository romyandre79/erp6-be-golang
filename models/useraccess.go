package models

import "time"

type Useraccess struct {
	Useraccessid int       `gorm:"column:useraccessid;primaryKey" json:"useraccessid"`
	Username     string    `gorm:"column:username" json:"username"`
	Realname     string    `gorm:"column:realname" json:"realname"`
	Password     string    `gorm:"column:password" json:"password"`
	Email        string    `gorm:"column:email" json:"email"`
	Phoneno      string    `gorm:"column:phoneno" json:"phoneno"`
	Languageid   int       `gorm:"column:languageid" json:"languageid"`
	Themeid      int       `gorm:"column:themeid" json:"themeid"`
	Isonline     int       `gorm:"column:isonline" json:"isonline"`
	Joindate     time.Time `gorm:"column:joindate" json:"joindate"`
	Authkey      string    `gorm:"column:authkey" json:"authkey"`
	Userphoto    string    `gorm:"column:userphoto" json:"userphoto"`
	Usertelegram string    `gorm:"column:usertelegram" json:"usertelegram"`
	Userwa       string    `gorm:"column:userwa" json:"userwa"`
	Wallpaper    string    `gorm:"column:wallpaper" json:"wallpaper"`
	Identityid   string    `gorm:"column:identityid" json:"identityid"`
	Signature    string    `gorm:"column:signature" json:"signature"`
	Recordstatus int       `gorm:"column:recordstatus" json:"recordstatus"`
	Lastlogin    time.Time `gorm:"column:lastlogin" json:"lastlogin"`
	Updatedate   time.Time `gorm:"column:updatedate" json:"updatedate"`
	Language     Language  `gorm:"foreignKey:languageid;references:languageid" json:"language"`
	Theme        Theme     `gorm:"foreignKey:themeid;references:themeid" json:"theme"`

	Usergroups []Usergroup `gorm:"foreignKey:Useraccessid;references:Useraccessid" json:"groups"`
}

func (Useraccess) TableName() string {
	return "useraccess"
}
