package common

import "erp6-be-golang/core/events"

func Initmaterialgroup () {
	events.Register("BeforeGet:materialgroup", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:materialgroup", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:materialgroup", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:materialgroup", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:materialgroup", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:materialgroup", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:materialgroup", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:materialgroup", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
