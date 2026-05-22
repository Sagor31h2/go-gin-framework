package models

import "time"

type NoteShare struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	NoteID         uint      `gorm:"not null;index" json:"note_id"`
	OwnerID        uint      `gorm:"not null" json:"owner_id"`
	SharedWithUser uint      `gorm:"not null" json:"shared_with_user"`
	CreatedAt      time.Time `json:"created_at"`
}
