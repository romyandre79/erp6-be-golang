package admin

import "erp6-be-golang/core/events"

func InitSnro() {
	events.Register("BeforeGet:snro", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:snro", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:snro", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:snro", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:snro", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:snro", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:snro", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:snro", func(data interface{}) error {
		return nil
	})
}
