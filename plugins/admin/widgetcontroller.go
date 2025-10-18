package admin

import "erp6-be-golang/core/events"

func InitWidget() {
	events.Register("BeforeGet:widget", func(data interface{}) error {
		return nil
	})
	events.Register("AfterGet:widget", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeSave:widget", func(data interface{}) error {
		return nil
	})
	events.Register("AfterSave:widget", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeUpdate:widget", func(data interface{}) error {
		return nil
	})
	events.Register("AfterUpdate:widget", func(data interface{}) error {
		return nil
	})
	events.Register("BeforeDelete:widget", func(data interface{}) error {
		return nil
	})
	events.Register("AfterDelete:widget", func(data interface{}) error {
		return nil
	})
}
