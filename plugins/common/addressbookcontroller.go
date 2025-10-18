package common

import "erp6-be-golang/core/events"

func Initaddressbook () {
	events.Register("BeforeGet:addressbook", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:addressbook", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:addressbook", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:addressbook", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:addressbook", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:addressbook", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:addressbook", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:addressbook", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
