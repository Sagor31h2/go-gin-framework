package models

type Note struct {
	ID     int    `gorm:"primaryKey" json:"id"`
	UserID uint   `json:"user_id"`
	Title  string `json:"title" binding:"required"`
	Status bool   `json:"status"`
}
