package models

import (
	"time"
)

type Report struct {
	ReportID int       `gorm:"column:reportid;primaryKey;autoIncrement" json:"reportid"`
	ReportName       string    `gorm:"column:reportname;not null" json:"reportname"`
	ReportDesc       string    `gorm:"column:reportdesc" json:"reportdesc"`
	ReportCategory   string    `gorm:"column:reportcategory" json:"reportcategory"`
	ReportType       string    `gorm:"column:reporttype" json:"reporttype"`                    // PDF, XLS, CSV
	TemplateJSON     string    `gorm:"column:templatejson;type:text" json:"templatejson"`      // Visual designer JSON
	JRXMLContent     string    `gorm:"column:jrxmlcontent;type:text" json:"jrxmlcontent"`      // JRXML format (for compatibility)
	Parameters       string    `gorm:"column:parameters;type:text" json:"parameters"`          // JSON array of parameters
	DataSource       string    `gorm:"column:datasource" json:"datasource"`                    // SQL query or workflow
	PageWidth        int       `gorm:"column:pagewidth;default:595" json:"pagewidth"`          // A4 width in points
	PageHeight       int       `gorm:"column:pageheight;default:842" json:"pageheight"`        // A4 height in points
	Orientation      string    `gorm:"column:orientation;default:portrait" json:"orientation"` // portrait or landscape
	ModuleID         int       `gorm:"column:moduleid" json:"moduleid"`
	RecordStatus     int8      `gorm:"column:recordstatus" json:"recordstatus"`
	UpdateDate       time.Time `gorm:"column:updatedate" json:"updatedate"`
}

func (Report) TableName() string {
	return "report"
}
