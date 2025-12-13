package models

import "time"

// ModuleWorkflow represents the moduleworkflow table
// Links workflows to modules for proper tracking and cleanup
type ModuleWorkflow struct {
	ModuleWorkflowID int       `gorm:"column:moduleworkflowid;primaryKey;autoIncrement" json:"moduleworkflowid"`
	ModuleID         int       `gorm:"column:moduleid;not null" json:"moduleid"`
	WorkflowID       int       `gorm:"column:workflowid;not null" json:"workflowid"`
	CreatedAt        time.Time `gorm:"column:createdat;not null;default:CURRENT_TIMESTAMP" json:"createdat"`

	// Foreign key relationships
	Module   *Modules  `gorm:"foreignKey:ModuleID;references:Moduleid" json:"module,omitempty"`
	Workflow *Workflow `gorm:"foreignKey:WorkflowID;references:Workflowid" json:"workflow,omitempty"`
}

// TableName specifies the table name for GORM
func (ModuleWorkflow) TableName() string {
	return "moduleworkflow"
}
