package common

import "erp6-be-golang/core/events"

func Initaddresstype () {
	events.Register("BeforeGet:addresstype", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:addresstype", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:addresstype", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:addresstype", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:addresstype", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:addresstype", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:addresstype", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:addresstype", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
