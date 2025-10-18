package common

import "erp6-be-golang/core/events"

func Initaddresscontact () {
	events.Register("BeforeGet:addresscontact", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:addresscontact", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:addresscontact", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:addresscontact", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:addresscontact", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:addresscontact", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:addresscontact", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:addresscontact", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
