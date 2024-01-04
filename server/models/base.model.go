package models

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=postgres user=%s password=%s dbname=%s port=5432 sslmode=prefer",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error! - ", err)
	}
	fmt.Println("Connected to DB")

	db.AutoMigrate(&User{}, &Msg{})
}
