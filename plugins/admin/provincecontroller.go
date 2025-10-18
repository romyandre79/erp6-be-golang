package admin

import "erp6-be-golang/core/events"

func InitProvince() {
	events.Register("BeforeGet:province", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:province", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:province", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:province", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:province", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:province", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:province", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:province", func(data interface{}) error {
		return nil
	})
}
