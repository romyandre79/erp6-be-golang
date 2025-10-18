package common

import "erp6-be-golang/core/events"

func Initownership () {
	events.Register("BeforeGet:ownership", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:ownership", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:ownership", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:ownership", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:ownership", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:ownership", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:ownership", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:ownership", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
