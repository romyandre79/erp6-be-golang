package blog

import "erp6-be-golang/core/events"

func Initpostcomment () {
	events.Register("BeforeGet:postcomment", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:postcomment", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:postcomment", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:postcomment", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:postcomment", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:postcomment", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:postcomment", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:postcomment", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
