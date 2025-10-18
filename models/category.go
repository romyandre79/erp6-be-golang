package models

import (
  "time"
)

type Category struct {
	Categoryid int `gorm:"column:categoryid;primaryKey" json:"categoryid"`
	Companyid int `gorm:"column:companyid" json:"companyid"`
	Title string `gorm:"column:title" json:"title"`
	Parentid *int `gorm:"column:parentid" json:"parentid"`
	Description string `gorm:"column:description" json:"description"`
	Slug string `gorm:"column:slug" json:"slug"`
	URL string `gorm:"column:url" json:"url"`
	Sort int8 `gorm:"column:sort" json:"sort"`
	Recordstatus int8 `gorm:"column:recordstatus" json:"recordstatus"`
	Updatedate time.Time `gorm:"column:updatedate" json:"updatedate"`
	Company *Company `gorm:"foreignKey:companyid;references:companyid"`
}

func (Category) TableName() string {
	return "category"
}
