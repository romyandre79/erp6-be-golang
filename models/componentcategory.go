package models

import (

)

type Componentcategory struct {
	Componentcategoryid int `gorm:"column:componentcategoryid;primaryKey" json:"componentcategoryid"`
	Categoryname string `gorm:"column:categoryname" json:"categoryname"`
}

func (Componentcategory) TableName() string {
	return "componentcategory"
}
