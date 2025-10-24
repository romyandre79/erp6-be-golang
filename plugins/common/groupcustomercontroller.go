package common

import "erp6-be-golang/core/events"

func Initgroupcustomer() {
	events.Register("BeforeGet:groupcustomer", func(data interface{}) error {
		// TODO: implement before get
		return nil
	})
	events.Register("AfterGet:groupcustomer", func(data interface{}) error {
		// TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:groupcustomer", func(data interface{}) error {
		// TODO: implement before save
		return nil
	})
	events.Register("AfterSave:groupcustomer", func(data interface{}) error {
		// TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:groupcustomer", func(data interface{}) error {
		// TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:groupcustomer", func(data interface{}) error {
		// TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:groupcustomer", func(data interface{}) error {
		// TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:groupcustomer", func(data interface{}) error {
		// TODO: implement after delete
		return nil
	})
}
