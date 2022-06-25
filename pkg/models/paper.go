package models

type Paper struct {
	ID        uint64  `gorm:"primary_key" json:"id"`
	Title     string  `gorm:"type:varchar(255)" json:"title"`
	Author    string  `gorm:"type:varchar(255)" json:"author"`
	Year      int64   `gorm:"type:int(4)" json:"year"`
	Abstract  string  `gorm:"type:text" json:"abstract"`
	File      string  `gorm:"type:varchar(255)" json:"file"`
	StudentID int64   `gorm:"type:int(11)" json:"-"`
	Student   Student `gorm:"foreignkey:StudentID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
}
