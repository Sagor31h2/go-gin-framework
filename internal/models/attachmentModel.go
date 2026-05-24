package models

import "time"

type Attachment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NoteID    int       `json:"note_id"`
	FileName  string    `json:"file_name"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}
