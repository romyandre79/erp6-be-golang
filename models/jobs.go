package models

import (
	"time"
)

type Jobs struct {
	JobsID       int       `gorm:"column:jobsid;primaryKey;autoIncrement:false" json:"jobsid"`
	JobName      string    `gorm:"column:jobname;type:varchar(50);not null;default:''" json:"jobname"`
	Description  string    `gorm:"column:description;type:varchar(100);not null;default:''" json:"description"`
	Version      string    `gorm:"column:version;type:varchar(100);not null;default:''" json:"version"`
	CreatedBy    string    `gorm:"column:createdby;type:varchar(100);not null;default:''" json:"createdby"`
	Flow         string    `gorm:"column:flow;type:varchar(100);not null;default:''" json:"flow"`
	Schedule     string    `gorm:"column:schedule;type:varchar(100);not null;default:''" json:"schedule"`
	IsRunning    int       `gorm:"column:isrunning;type:tinyint(4);not null" json:"is_running"`
	LastRunning  time.Time `gorm:"column:lastrunning;type:timestamp;not null" json:"last_running"`
	RecordStatus int       `gorm:"column:recordstatus;type:tinyint(4);not null;default:0" json:"record_status"`
}

func (Jobs) TableName() string {
	return "jobs"
}
