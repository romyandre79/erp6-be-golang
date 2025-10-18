package admin

import "erp6-be-golang/core/events"

func InitMenuaccess() {
	events.Register("BeforeGet:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:menuaccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:menuaccess", func(data interface{}) error {
		return nil
	})
}
