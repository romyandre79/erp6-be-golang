package models

type DbobjectRelation struct {
	ID           int    `gorm:"column:dbobjectrelationid;primaryKey;autoIncrement" json:"id"`
	FromTableID  int    `gorm:"column:fromtableid" json:"from_table_id"`
	FromColIndex int    `gorm:"column:fromcolindex" json:"from_col_index"`
	FromColName  string `gorm:"column:fromcolname;type:varchar(255)" json:"from_col_name"`
	ToTableID    int    `gorm:"column:totableid" json:"to_table_id"`
	ToColIndex   int    `gorm:"column:tocolindex" json:"to_col_index"`
	ToColName    string `gorm:"column:tocolname;type:varchar(255)" json:"to_col_name"`
	Path         string `gorm:"column:path;type:text" json:"path"` // SVG path data or routing info
}

func (DbobjectRelation) TableName() string {
	return "dbobjectrelation"
}
