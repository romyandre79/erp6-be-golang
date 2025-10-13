package admin

import "time"

type UserAccess struct {
	UserAccessID int       `gorm:"column:useraccessid;primaryKey;autoIncrement" json:"useraccessid"`
	Username     string    `gorm:"column:username;unique;not null" json:"username"`
	RealName     string    `gorm:"column:realname" json:"realname"`
	Password     string    `gorm:"column:password" json:"-"`
	Email        string    `gorm:"column:email" json:"email"`
	PhoneNo      string    `gorm:"column:phoneno" json:"phoneno"`
	LanguageID   int       `gorm:"column:languageid;not null" json:"languageid"`
	ThemeID      int       `gorm:"column:themeid;not null" json:"themeid"`
	IsOnline     bool      `gorm:"column:isonline;default:0" json:"isonline"`
	JoinDate     time.Time `gorm:"column:joindate;autoCreateTime" json:"joindate"`
	AuthKey      string    `gorm:"column:authkey" json:"authkey"`
	UserPhoto    string    `gorm:"column:userphoto;default:man.png" json:"userphoto"`
	UserTelegram string    `gorm:"column:usertelegram" json:"usertelegram"`
	UserWA       string    `gorm:"column:userwa" json:"userwa"`
	Wallpaper    string    `gorm:"column:wallpaper;default:Andromeda-Galaxy.jpg" json:"wallpaper"`
	IdentityID   string    `gorm:"column:identityid" json:"identityid"`
	Signature    string    `gorm:"column:signature" json:"signature"`
	RecordStatus int       `gorm:"column:recordstatus;default:0" json:"recordstatus"`
	LastLogin    time.Time `gorm:"column:lastlogin;autoUpdateTime" json:"lastlogin"`
	UpdateDate   time.Time `gorm:"column:updatedate;autoUpdateTime" json:"updatedate"`

	// Relasi ke UserGroup
	Groups []UserGroup `gorm:"foreignKey:UserAccessID;references:UserAccessID" json:"groups"`
}

func (UserAccess) TableName() string {
	return "useraccess"
}
