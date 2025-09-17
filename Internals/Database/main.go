package database

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to DB:", err) // <-- show actual error
		panic("failed to connect db")
	}

	fmt.Println("DB connected successfully")
	return db
}
