package common

import "erp6-be-golang/core/events"

func Initproduct () {
	events.Register("BeforeGet:product", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:product", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:product", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:product", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:product", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:product", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:product", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:product", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
