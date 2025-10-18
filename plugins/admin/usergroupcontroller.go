package admin

import "erp6-be-golang/core/events"

func InitUsergroup() {
	events.Register("BeforeGet:usergroup", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:usergroup", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:usergroup", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:usergroup", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:usergroup", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:usergroup", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:usergroup", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:usergroup", func(data interface{}) error {
		return nil
	})
}
