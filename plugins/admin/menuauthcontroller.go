package admin

import "erp6-be-golang/core/events"

func InitMenuauth() {
	events.Register("BeforeGet:menuauth", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:menuauth", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:menuauth", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:menuauth", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:menuauth", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:menuauth", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:menuauth", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:menuauth", func(data interface{}) error {
		return nil
	})
}
