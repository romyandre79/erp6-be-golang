package admin

import "erp6-be-golang/core/events"

func InitCity() {
	events.Register("BeforeGet:city", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:city", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:city", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:city", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:city", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:city", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:city", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:city", func(data interface{}) error {
		return nil
	})
}
