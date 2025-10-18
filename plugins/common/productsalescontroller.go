package common

import "erp6-be-golang/core/events"

func Initproductsales() {
	events.Register("BeforeGet:productsales", func(data interface{}) error {
		// TODO: implement before get
		return nil
	})
	events.Register("AfterGet:productsales", func(data interface{}) error {
		// TODO: implement after get
		return nil
	})
	events.Register("BeforeSave:productsales", func(data interface{}) error {
		// TODO: implement before save
		return nil
	})
	events.Register("AfterSave:productsales", func(data interface{}) error {
		// TODO: implement after save
		return nil
	})
	events.Register("BeforeUpdate:productsales", func(data interface{}) error {
		// TODO: implement before update
		return nil
	})
	events.Register("AfterUpdate:productsales", func(data interface{}) error {
		// TODO: implement after update
		return nil
	})
	events.Register("BeforeDelete:productsales", func(data interface{}) error {
		// TODO: implement before delete
		return nil
	})
	events.Register("AfterDelete:productsales", func(data interface{}) error {
		// TODO: implement after delete
		return nil
	})
}
