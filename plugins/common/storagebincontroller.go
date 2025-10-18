package common

import "erp6-be-golang/core/events"

func Initstoragebin () {
	events.Register("BeforeGet:storagebin", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:storagebin", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:storagebin", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:storagebin", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:storagebin", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:storagebin", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:storagebin", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:storagebin", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
