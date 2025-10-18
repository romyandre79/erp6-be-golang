package models

import (

)

type Componentdetail struct {
	Componentdetailid int `gorm:"column:componentdetailid;primaryKey" json:"componentdetailid"`
	Componentid int `gorm:"column:componentid" json:"componentid"`
	Detailtype *string `gorm:"column:detailtype" json:"detailtype"`
	Lable *string `gorm:"column:lable" json:"lable"`
	Inputtype *string `gorm:"column:inputtype" json:"inputtype"`
	Inputname *string `gorm:"column:inputname" json:"inputname"`
	Inputdesc *string `gorm:"column:inputdesc" json:"inputdesc"`
	Order *int `gorm:"column:order" json:"order"`
	Datasourcetype *string `gorm:"column:datasourcetype" json:"datasourcetype"`
	Datasource *string `gorm:"column:datasource" json:"datasource"`
	Datasourceidfield *string `gorm:"column:datasourceidfield" json:"datasourceidfield"`
	Datasourcenamefield *string `gorm:"column:datasourcenamefield" json:"datasourcenamefield"`
	Component Component `gorm:"foreignKey:componentid;references:componentid"`
}

func (Componentdetail) TableName() string {
	return "componentdetail"
}
