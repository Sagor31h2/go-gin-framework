package models

type Note struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	Title  string `json:"title" binding:"required"`
	Status bool   `json:"status"`
}
