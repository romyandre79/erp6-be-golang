package common

import "erp6-be-golang/core/events"

func Initaddress () {
	events.Register("BeforeGet:address", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:address", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:address", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:address", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:address", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:address", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:address", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:address", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
