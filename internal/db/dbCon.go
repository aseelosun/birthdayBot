package db

import (
	"birthdayBot/internal/configs"
	"database/sql"
	"fmt"
	"log"
)

func ConnectToDb(conf configs.Configuration) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.DbConfig.Host, conf.DbConfig.Port, conf.DbConfig.User, conf.DbConfig.Password, conf.DbConfig.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	fmt.Printf("Postgres Connected!\n")

	return db, err
}
