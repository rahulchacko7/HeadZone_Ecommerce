package domain

type Category struct {
	Id       uint   `json:"id" gorm:"unique; not null"`
	Category string `json:"category" gorm:"unique;not null"`
}
