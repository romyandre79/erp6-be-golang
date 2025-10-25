package response

import (
	"time"
)

type Post struct {
	Postid       int       `json:"postid"`
	Companyid    int       `json:"companyid"`
	Postpic      string    `json:"postpic"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Metatag      string    `json:"metatag"`
	Slug         string    `json:"slug"`
	Postupdate   time.Time `json:"postupdate"`
	Created      time.Time `json:"created"`
	Recordstatus int8      `json:"recordstatus"`
	Author       Author    `json:"author"`
	Category     string    `json:"category"`
}

type Author struct {
	Useraccessid int    `json:"useraccessid"`
	Realname     string `json:"realname"`
}

type Category struct {
	Categoryid   int    `json:"categoryid"`
	Categoryname string `json:"categoryname"`
}
