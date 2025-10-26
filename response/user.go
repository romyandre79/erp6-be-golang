package response

type UserResponse struct {
	UserAccessID uint                `json:"useraccessid"`
	Username     string              `json:"username"`
	RealName     string              `json:"realname"`
	Email        string              `json:"email"`
	Phoneno      string              `json:"phoneno"`
	Photo        string              `json:"photo"`
	Language     string              `json:"language"`
	Groups       []GroupResponse     `json:"groups"`
	Menus        []GroupMenuResponse `json:"menus"`
}

type GroupMenuResponse struct {
	MenuAccessID int    `json:"menuaccessid"`
	MenuName     string `json:"menuname"`
	MenuCode     string `json:"menucode"`
	ParentId     *int   `json:"parentid"`
	IsRead       int    `json:"isread"`
	IsWrite      int    `json:"iswrite"`
	IsPost       int    `json:"ispost"`
	IsReject     int    `json:"isreject"`
	IsUpload     int    `json:"isupload"`
	IsDownload   int    `json:"isdownload"`
	IsPurge      int    `json:"ispurge"`
}

type GroupResponse struct {
	GroupAccessID int    `json:"groupaccessid"`
	GroupName     string `json:"groupname"`
	Description   string `json:"description"`
}
