package admin

import "erp6-be-golang/core/events"

func InitGroupmenuauth() {
	events.Register("BeforeGet:groupmenuauth", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:groupmenuauth", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:groupmenuauth", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:groupmenuauth", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:groupmenuauth", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:groupmenuauth", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:groupmenuauth", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:groupmenuauth", func(data interface{}) error {
		return nil
	})
}
