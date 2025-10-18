package admin

import "erp6-be-golang/core/events"

func InitGroupmenu() {
	events.Register("BeforeGet:groupmenu", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:groupmenu", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:groupmenu", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:groupmenu", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:groupmenu", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:groupmenu", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:groupmenu", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:groupmenu", func(data interface{}) error {
		return nil
	})
}
