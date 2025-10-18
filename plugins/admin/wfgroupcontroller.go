package admin

import "erp6-be-golang/core/events"

func InitWfgroup() {
	events.Register("BeforeGet:wfgroup", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:wfgroup", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:wfgroup", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:wfgroup", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:wfgroup", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:wfgroup", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:wfgroup", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:wfgroup", func(data interface{}) error {
		return nil
	})
}
