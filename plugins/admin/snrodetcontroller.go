package admin

import "erp6-be-golang/core/events"

func InitSnrodet() {
	events.Register("BeforeGet:snrodet", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:snrodet", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:snrodet", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:snrodet", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:snrodet", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:snrodet", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:snrodet", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:snrodet", func(data interface{}) error {
		return nil
	})
}
