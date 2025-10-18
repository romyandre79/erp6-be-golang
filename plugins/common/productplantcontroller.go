package common

import "erp6-be-golang/core/events"

func Initproductplant () {
	events.Register("BeforeGet:productplant", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:productplant", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:productplant", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:productplant", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:productplant", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:productplant", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:productplant", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:productplant", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
