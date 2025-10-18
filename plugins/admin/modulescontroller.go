package admin

import "erp6-be-golang/core/events"

func InitModules() {
	events.Register("BeforeGet:modules", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:modules", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:modules", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:modules", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:modules", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:modules", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:modules", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:modules", func(data interface{}) error {
		return nil
	})
}
