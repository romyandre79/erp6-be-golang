package response

type GroupResponse struct {
	GroupAccessID int    `gorm:"column:groupaccessid;primaryKey;autoIncrement"`
	GroupName     string `gorm:"column:groupname;size:50;unique;not null"`
	Description   string `gorm:"column:description;size:100;not null"`
}
