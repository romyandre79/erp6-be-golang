package models

type Reportout struct {
	Reportoutid   int         `gorm:"column:reportoutid;primaryKey" json:"reportoutid"`
	Addressbookid int         `gorm:"column:addressbookid" json:"addressbookid"`
	Bulan         int         `gorm:"column:bulan" json:"bulan"`
	Tahun         int         `gorm:"column:tahun" json:"tahun"`
	S1            string      `gorm:"column:s1" json:"s1"`
	D1            string      `gorm:"column:d1" json:"d1"`
	S2            string      `gorm:"column:s2" json:"s2"`
	D2            string      `gorm:"column:d2" json:"d2"`
	S3            string      `gorm:"column:s3" json:"s3"`
	D3            string      `gorm:"column:d3" json:"d3"`
	S4            string      `gorm:"column:s4" json:"s4"`
	D4            string      `gorm:"column:d4" json:"d4"`
	S5            string      `gorm:"column:s5" json:"s5"`
	D5            string      `gorm:"column:d5" json:"d5"`
	S6            string      `gorm:"column:s6" json:"s6"`
	D6            string      `gorm:"column:d6" json:"d6"`
	S7            string      `gorm:"column:s7" json:"s7"`
	D7            string      `gorm:"column:d7" json:"d7"`
	S8            string      `gorm:"column:s8" json:"s8"`
	D8            string      `gorm:"column:d8" json:"d8"`
	S9            string      `gorm:"column:s9" json:"s9"`
	D9            string      `gorm:"column:d9" json:"d9"`
	S10           string      `gorm:"column:s10" json:"s10"`
	D10           string      `gorm:"column:d10" json:"d10"`
	S11           string      `gorm:"column:s11" json:"s11"`
	D11           string      `gorm:"column:d11" json:"d11"`
	S12           string      `gorm:"column:s12" json:"s12"`
	D12           string      `gorm:"column:d12" json:"d12"`
	S13           string      `gorm:"column:s13" json:"s13"`
	D13           string      `gorm:"column:d13" json:"d13"`
	S14           string      `gorm:"column:s14" json:"s14"`
	D14           string      `gorm:"column:d14" json:"d14"`
	S15           string      `gorm:"column:s15" json:"s15"`
	D15           string      `gorm:"column:d15" json:"d15"`
	S16           string      `gorm:"column:s16" json:"s16"`
	D16           string      `gorm:"column:d16" json:"d16"`
	S17           string      `gorm:"column:s17" json:"s17"`
	D17           string      `gorm:"column:d17" json:"d17"`
	S18           string      `gorm:"column:s18" json:"s18"`
	D18           string      `gorm:"column:d18" json:"d18"`
	S19           string      `gorm:"column:s19" json:"s19"`
	D19           string      `gorm:"column:d19" json:"d19"`
	S20           string      `gorm:"column:s20" json:"s20"`
	D20           string      `gorm:"column:d20" json:"d20"`
	S21           string      `gorm:"column:s21" json:"s21"`
	D21           string      `gorm:"column:d21" json:"d21"`
	S22           string      `gorm:"column:s22" json:"s22"`
	D22           string      `gorm:"column:d22" json:"d22"`
	S23           string      `gorm:"column:s23" json:"s23"`
	D23           string      `gorm:"column:d23" json:"d23"`
	S24           string      `gorm:"column:s24" json:"s24"`
	D24           string      `gorm:"column:d24" json:"d24"`
	S25           string      `gorm:"column:s25" json:"s25"`
	D25           string      `gorm:"column:d25" json:"d25"`
	S26           string      `gorm:"column:s26" json:"s26"`
	D26           string      `gorm:"column:d26" json:"d26"`
	S27           string      `gorm:"column:s27" json:"s27"`
	D27           string      `gorm:"column:d27" json:"d27"`
	S28           string      `gorm:"column:s28" json:"s28"`
	D28           string      `gorm:"column:d28" json:"d28"`
	S29           string      `gorm:"column:s29" json:"s29"`
	D29           string      `gorm:"column:d29" json:"d29"`
	S30           string      `gorm:"column:s30" json:"s30"`
	D30           string      `gorm:"column:d30" json:"d30"`
	S31           string      `gorm:"column:s31" json:"s31"`
	D31           string      `gorm:"column:d31" json:"d31"`
	Oldnik        string      `gorm:"column:oldnik" json:"oldnik"`
	Addressbook   Addressbook `gorm:"foreignKey:addressbookid;references:addressbookid"`
}

func (Reportout) TableName() string {
	return "reportout"
}
