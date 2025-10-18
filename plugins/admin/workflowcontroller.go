package admin

import "erp6-be-golang/core/events"

func InitWorkflow() {
	events.Register("BeforeGet:workflow", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:workflow", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:workflow", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:workflow", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:workflow", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:workflow", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:workflow", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:workflow", func(data interface{}) error {
		return nil
	})
}
