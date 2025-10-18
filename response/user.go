package response

type UserResponse struct {
	UserAccessID uint                `json:"useraccessid"`
	Username     string              `json:"username"`
	RealName     string              `json:"realname"`
	Email        string              `json:"email"`
	Phoneno      string              `json:"phoneno"`
	Language     string              `json:"language"`
	Groups       []GroupResponse     `json:"groups"`
	Menus        []GroupMenuResponse `json:"menus"`
}
