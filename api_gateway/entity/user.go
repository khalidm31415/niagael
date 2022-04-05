package entity

type User struct {
	Base
	Username     string `gorm:"type:VARCHAR; NOT NULL; unique"`
	DisplayName  string `gorm:"type:VARCHAR; NOT NULL"`
	PasswordHash string `gorm:"type:TEXT; NOT NULL" json:"-"`
}
