package admin

import "erp6-be-golang/core/events"

func InitUseraccess() {
	events.Register("BeforeGet:useraccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:useraccess", func(data interface{}) error {
		return nil
	})
	events.Register("BseforeSave:useraccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:useraccess", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:useraccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:useraccess", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:useraccess", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:useraccess", func(data interface{}) error {
		return nil
	})
}
