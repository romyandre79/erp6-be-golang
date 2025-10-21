package models

type Component struct {
	Componentid         int               `gorm:"column:componentid;primaryKey" json:"componentid"`
	Componentname       string            `gorm:"column:componentname" json:"componentname"`
	Componenttitle      string            `gorm:"column:componenttitle" json:"componenttitle"`
	Componentcategoryid int               `gorm:"column:componentcategoryid" json:"componentcategoryid"`
	Componentclass      string            `gorm:"column:componentclass" json:"componentclass"`
	Input               int8              `gorm:"column:input" json:"input"`
	Output              int8              `gorm:"column:output" json:"output"`
	Componentcategory   Componentcategory `gorm:"foreignKey:componentcategoryid;references:componentcategoryid"`

	Details []Componentdetail `gorm:"foreignKey:Componentid;references:Componentid" json:"details"`
}

func (Component) TableName() string {
	return "component"
}
