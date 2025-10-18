package common

import "erp6-be-golang/core/events"

func Initidentitytype () {
	events.Register("BeforeGet:identitytype", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:identitytype", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:identitytype", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:identitytype", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:identitytype", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:identitytype", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:identitytype", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:identitytype", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
