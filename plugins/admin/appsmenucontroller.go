package admin

import "erp6-be-golang/core/events"

func InitAppsMenu() {
	events.Register("BeforeGet:appsmenu", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:appsmenu", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:appsmenu", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:appsmenu", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:appsmenu", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:appsmenu", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:appsmenu", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:appsmenu", func(data interface{}) error {
		return nil
	})
}
