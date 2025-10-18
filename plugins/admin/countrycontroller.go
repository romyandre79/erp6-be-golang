package admin

import "erp6-be-golang/core/events"

func InitCountry() {
	events.Register("BeforeGet:country", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:country", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:country", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:country", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:country", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:country", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:country", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:country", func(data interface{}) error {
		return nil
	})
}
