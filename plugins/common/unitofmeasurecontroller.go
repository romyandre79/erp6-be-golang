package common

import "erp6-be-golang/core/events"

func Initunitofmeasure () {
	events.Register("BeforeGet:unitofmeasure", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:unitofmeasure", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:unitofmeasure", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:unitofmeasure", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:unitofmeasure", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:unitofmeasure", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:unitofmeasure", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:unitofmeasure", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
