package admin

import "erp6-be-golang/core/events"

func InitApps() {
	events.Register("BeforeGet:apps", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:apps", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:apps", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:apps", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:apps", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:apps", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:apps", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:apps", func(data interface{}) error {
		return nil
	})
}
