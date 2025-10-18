package admin

import "erp6-be-golang/core/events"

func InitTranslog() {
	events.Register("BeforeGet:translog", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:translog", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:translog", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:translog", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:translog", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:translog", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:translog", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:translog", func(data interface{}) error {
		return nil
	})
}
