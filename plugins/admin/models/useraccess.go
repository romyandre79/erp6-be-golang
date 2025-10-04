package admin

import "time"

type UserAccess struct {
	UserAccessID int       `gorm:"column:useraccessid;primaryKey" json:"useraccessid"`
	Username     string    `gorm:"column:username" json:"username"`
	RealName     string    `gorm:"column:realname" json:"realname"`
	Password     string    `gorm:"column:password" json:"-"`
	Email        string    `gorm:"column:email" json:"email"`
	PhoneNo      string    `gorm:"column:phoneno" json:"phoneno"`
	LanguageID   int       `gorm:"column:languageid" json:"languageid"`
	ThemeID      int       `gorm:"column:themeid" json:"themeid"`
	IsOnline     int       `gorm:"column:isonline" json:"isonline"`
	JoinDate     time.Time `gorm:"column:joindate" json:"joindate"`
	AuthKey      string    `gorm:"column:authkey" json:"authkey"`
	UserPhoto    string    `gorm:"column:userphoto" json:"userphoto"`
	UserTelegram string    `gorm:"column:usertelegram" json:"usertelegram"`
	UserWA       string    `gorm:"column:userwa" json:"userwa"`
	Wallpaper    string    `gorm:"column:wallpaper" json:"wallpaper"`
	IdentityID   string    `gorm:"column:identityid" json:"identityid"`
	Signature    string    `gorm:"column:signature" json:"signature"`
	RecordStatus int       `gorm:"column:recordstatus" json:"recordstatus"`
	LastLogin    time.Time `gorm:"column:lastlogin" json:"lastlogin"`
	UpdateDate   time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (UserAccess) TableName() string {
	return "useraccess"
}
