package models

type Student struct {
	ID       uint64 `gorm:"primary_key" json:"id"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	NRP      string `gorm:"type:varchar(255)" json:"nrp"`
	Email    string `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password string `gorm:"->;<-;not null" json:"-"`
	Token    string `gorm:"-" json:"token,omitempty"`
}
