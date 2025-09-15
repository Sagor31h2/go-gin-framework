package services

type NotesService struct{}

type Note struct {
	Id   int
	Name string
}

func (n *NotesService) GetNotes() []Note {
	return []Note{
		{Id: 1, Name: "Note 1"},
		{Id: 2, Name: "Note "},
		{Id: 3, Name: "Note 3"},
		{Id: 1, Name: "Note 1"},
		{Id: 1, Name: "Note 1"},
		{Id: 1, Name: "Note 1"},
		{Id: 1, Name: "Note 1"},
	}
}
