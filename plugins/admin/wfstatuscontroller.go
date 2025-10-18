package admin

import "erp6-be-golang/core/events"

func InitWfstatus() {
	events.Register("BeforeGet:wfstatus", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:wfstatus", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:wfstatus", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:wfstatus", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:wfstatus", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:wfstatus", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:wfstatus", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:wfstatus", func(data interface{}) error {
		return nil
	})
}
