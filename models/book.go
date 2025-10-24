package models

import (
	"time"
)

type Book struct {
	Bookid          int           `gorm:"column:bookid;primaryKey" json:"bookid"`
	Companyid       int           `gorm:"column:companyid" json:"companyid"`
	Bookcode        string        `gorm:"column:bookcode" json:"bookcode"`
	Bookcategoryid  int           `gorm:"column:bookcategoryid" json:"bookcategoryid"`
	Bookmediatypeid int           `gorm:"column:bookmediatypeid" json:"bookmediatypeid"`
	Bookauthors     string        `gorm:"column:bookauthors" json:"bookauthors"`
	Bookeditor      string        `gorm:"column:bookeditor" json:"bookeditor"`
	Booktitle       string        `gorm:"column:booktitle" json:"booktitle"`
	Edition         string        `gorm:"column:edition" json:"edition"`
	Specdetailinfo  string        `gorm:"column:specdetailinfo" json:"specdetailinfo"`
	Isbn            string        `gorm:"column:isbn" json:"isbn"`
	Addressbookid   int           `gorm:"column:addressbookid" json:"addressbookid"`
	Publishingyear  int           `gorm:"column:publishingyear" json:"publishingyear"`
	Languageid      int           `gorm:"column:languageid" json:"languageid"`
	Abstractnotes   string        `gorm:"column:abstractnotes" json:"abstractnotes"`
	Imageurl        string        `gorm:"column:imageurl" json:"imageurl"`
	Bookrent        int           `gorm:"column:bookrent" json:"bookrent"`
	Bookcopy        int           `gorm:"column:bookcopy" json:"bookcopy"`
	Tableofcontent  string        `gorm:"column:tableofcontent" json:"tableofcontent"`
	Bookprice       int           `gorm:"column:bookprice" json:"bookprice"`
	Currencyid      int           `gorm:"column:currencyid" json:"currencyid"`
	Bookrentprice   int           `gorm:"column:bookrentprice" json:"bookrentprice"`
	Bookreturnday   int           `gorm:"column:bookreturnday" json:"bookreturnday"`
	Recordstatus    int8          `gorm:"column:recordstatus" json:"recordstatus"`
	Createddate     time.Time     `gorm:"column:createddate" json:"createddate"`
	Updatedate      time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Addressbook     Addressbook   `gorm:"foreignKey:addressbookid;references:addressbookid"`
	Bookcategory    Bookcategory  `gorm:"foreignKey:bookcategoryid;references:bookcategoryid"`
	Company         Company       `gorm:"foreignKey:companyid;references:companyid"`
	Language        Language      `gorm:"foreignKey:languageid;references:languageid"`
	Bookmediatype   Bookmediatype `gorm:"foreignKey:bookmediatypeid;references:bookmediatypeid"`
}

func (Book) TableName() string {
	return "book"
}
