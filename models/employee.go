package models

import (

)

type Employee struct {
}

func (Employee) TableName() string {
	return "employee"
}
