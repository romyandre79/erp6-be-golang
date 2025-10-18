package common

import "erp6-be-golang/core/events"

func Initsloc () {
	events.Register("BeforeGet:sloc", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:sloc", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:sloc", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:sloc", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:sloc", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:sloc", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:sloc", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:sloc", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
