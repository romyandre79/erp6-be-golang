package common

import "erp6-be-golang/core/events"

func Initmaterialstatus () {
	events.Register("BeforeGet:materialstatus", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:materialstatus", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:materialstatus", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:materialstatus", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:materialstatus", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:materialstatus", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:materialstatus", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:materialstatus", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
