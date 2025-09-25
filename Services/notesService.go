package services

import (
	models "gin-test/Internals/Models"

	"gorm.io/gorm"
)

type NotesService struct {
	db *gorm.DB
}

func (n *NotesService) InitNotesService(dataBase *gorm.DB) {
	n.db = dataBase
	n.db.AutoMigrate(&models.Note{})
}

func (n *NotesService) GetNotes() []models.Note {
	return []models.Note{
		{Id: 1, Title: "Note 1"},
		{Id: 2, Title: "Note "},
		{Id: 3, Title: "Note 3"},
		{Id: 1, Title: "Note 1"},
		{Id: 1, Title: "Note 1"},
		{Id: 1, Title: "Note 1"},
		{Id: 1, Title: "Note 1"},
	}
}

func (n *NotesService) CreateNotes() models.Note {
	data := models.Note{
		Title:  "test note 1",
		Status: true,
	}
	result := n.db.Create(&data)
	if result.Error != nil {
		println(result.Error.Error())
	}
	return data
}
