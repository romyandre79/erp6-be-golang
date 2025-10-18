package blog

import "erp6-be-golang/core/events"

func Initcategory () {
	events.Register("BeforeGet:category", func(data interface{}) error {
	    // TODO: implement before get
		return nil
	})
	events.Register("AfterGet:category", func(data interface{}) error {
	    // TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:category", func(data interface{}) error {
	    // TODO: implement before save
		return nil
	})
	events.Register("AfterSave:category", func(data interface{}) error {
	    // TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:category", func(data interface{}) error {
	    // TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:category", func(data interface{}) error {
	    // TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:category", func(data interface{}) error {
	    // TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:category", func(data interface{}) error {
	    // TODO: implement after delete
		return nil
	})
}
