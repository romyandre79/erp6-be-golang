package common

import "erp6-be-golang/core/events"

func Initromawi () {
	events.Register("BeforeGet:romawi", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:romawi", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:romawi", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:romawi", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:romawi", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:romawi", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:romawi", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:romawi", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
