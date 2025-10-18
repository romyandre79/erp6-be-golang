package admin

import "erp6-be-golang/core/events"

func InitCurrency() {
	events.Register("BeforeGet:currency", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:currency", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:currency", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:currency", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:currency", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:currency", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:currency", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:currency", func(data interface{}) error {
		return nil
	})
}
