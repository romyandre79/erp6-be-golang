package admin

import "erp6-be-golang/core/events"

func InitLanguage() {
	events.Register("BeforeGet:language", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:language", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:language", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:language", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:language", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:language", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:language", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:language", func(data interface{}) error {
		return nil
	})
}
