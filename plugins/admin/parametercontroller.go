package admin

import "erp6-be-golang/core/events"

func InitParameter() {
	events.Register("BeforeGet:parameter", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:parameter", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:parameter", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:parameter", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:parameter", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:parameter", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:parameter", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:parameter", func(data interface{}) error {
		return nil
	})
}
