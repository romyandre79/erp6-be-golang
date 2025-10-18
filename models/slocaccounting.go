package models

import (
	"time"
)

type Slocaccounting struct {
	Slocaccid              int           `gorm:"column:slocaccid;primaryKey" json:"slocaccid"`
	Slocid                 int           `gorm:"column:slocid" json:"slocid"`
	Materialgroupid        int           `gorm:"column:materialgroupid" json:"materialgroupid"`
	Accpembelianid         *int          `gorm:"column:accpembelian" json:"accpembelian"`
	Accpersediaanid        *int          `gorm:"column:accpersediaan" json:"accpersediaan"`
	Accreturpembelianid    *int          `gorm:"column:accreturpembelian" json:"accreturpembelian"`
	Accdiscpembelianid     *int          `gorm:"column:accdiscpembelian" json:"accdiscpembelian"`
	Accbiayapembelianid    *int          `gorm:"column:accbiayapembelian" json:"accbiayapembelian"`
	Accexpedisipembelianid *int          `gorm:"column:accexpedisipembelian" json:"accexpedisipembelian"`
	Accpenjualanid         *int          `gorm:"column:accpenjualan" json:"accpenjualan"`
	Accreturpenjualanid    *int          `gorm:"column:accreturpenjualan" json:"accreturpenjualan"`
	Accdiscpenjualanid     *int          `gorm:"column:accdiscpenjualan" json:"accdiscpenjualan"`
	Accbiayapenjualanid    *int          `gorm:"column:accbiayapenjualan" json:"accbiayapenjualan"`
	Accexpedisipenjualanid *int          `gorm:"column:accexpedisipenjualan" json:"accexpedisipenjualan"`
	Hppid                  *int          `gorm:"column:hpp" json:"hpp"`
	Accupahlemburid        *int          `gorm:"column:accupahlembur" json:"accupahlembur"`
	Fohid                  *int          `gorm:"column:foh" json:"foh"`
	Accakumid              *int          `gorm:"column:accakum" json:"accakum"`
	Accbiayasusutid        *int          `gorm:"column:accbiayasusut" json:"accbiayasusut"`
	Acckorsusutid          *int          `gorm:"column:acckorsusut" json:"acckorsusut"`
	Updatedate             time.Time     `gorm:"column:updatedate" json:"updatedate"`
	Accdiscpembelian       Account       `gorm:"foreignKey:accdiscpembelian;references:accountid"`
	Accpembelian           Account       `gorm:"foreignKey:accpembelian;references:accountid"`
	Accpersediaan          Account       `gorm:"foreignKey:accpersediaan;references:accountid"`
	Accreturpembelian      Account       `gorm:"foreignKey:accreturpembelian;references:accountid"`
	Materialgroup          Materialgroup `gorm:"foreignKey:materialgroupid;references:materialgroupid"`
	Sloc                   Sloc          `gorm:"foreignKey:slocid;references:slocid"`
}

func (Slocaccounting) TableName() string {
	return "slocaccounting"
}
