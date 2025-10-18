package admin

import "erp6-be-golang/core/events"

func InitTheme() {
	events.Register("BeforeGet:theme", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:theme", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:theme", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:theme", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:theme", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:theme", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:theme", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:theme", func(data interface{}) error {
		return nil
	})
}
