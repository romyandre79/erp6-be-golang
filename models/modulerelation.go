package models

type Modulerelation struct {
	Modulerelationid int     `gorm:"column:modulerelationid;primaryKey" json:"modulerelationid"`
	Moduleid         int     `gorm:"column:moduleid" json:"moduleid"`
	Relationid       int     `gorm:"column:relationid" json:"relationid"`
	Module           Modules `gorm:"foreignKey:moduleid;references:moduleid"`
	RelatedModule    Modules `gorm:"foreignKey:relationid;references:moduleid"`
}

func (Modulerelation) TableName() string {
	return "modulerelation"
}
