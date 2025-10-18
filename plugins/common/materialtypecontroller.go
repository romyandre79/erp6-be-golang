package common

import "erp6-be-golang/core/events"

func Initmaterialtype () {
	events.Register("BeforeGet:materialtype", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:materialtype", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:materialtype", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:materialtype", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:materialtype", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:materialtype", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:materialtype", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:materialtype", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
