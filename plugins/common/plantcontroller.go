package common

import "erp6-be-golang/core/events"

func Initplant () {
	events.Register("BeforeGet:plant", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:plant", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:plant", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:plant", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:plant", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:plant", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:plant", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:plant", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
