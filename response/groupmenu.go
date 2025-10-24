package response

type GroupMenuResponse struct {
	MenuAccessID int    `gorm:"column:menuaccessid;not null" json:"menuaccessid"`
	MenuName     string `gorm:"column:menuname;not null" json:"menuname"`
	MenuCode     string `gorm:"column:menucode;not null" json:"menucode"`
	ParentId     *int   `gorm:"column:parentid;not null" json:"parentid"`
	IsRead       int    `gorm:"column:isread;default:1" json:"isread"`
	IsWrite      int    `gorm:"column:iswrite;default:1" json:"iswrite"`
	IsPost       int    `gorm:"column:ispost;default:1" json:"ispost"`
	IsReject     int    `gorm:"column:isreject;default:1" json:"isreject"`
	IsUpload     int    `gorm:"column:isupload;default:1" json:"isupload"`
	IsDownload   int    `gorm:"column:isdownload;default:1" json:"isdownload"`
	IsPurge      int    `gorm:"column:ispurge;default:1" json:"ispurge"`
}
