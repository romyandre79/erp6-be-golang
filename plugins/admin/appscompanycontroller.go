package admin

import "erp6-be-golang/core/events"

func InitAppsCompany() {
	events.Register("BeforeGet:appscompany", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:appscompany", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:appscompany", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:appscompany", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:appscompany", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:appscompany", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:appscompany", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:appscompany", func(data interface{}) error {
		return nil
	})
}
