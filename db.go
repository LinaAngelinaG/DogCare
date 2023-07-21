package main

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "3846936720"
	dbname   = "dogcare"
)

type pet struct {
	userName  string
	petNumber int
	petName   string
	birthday  string
	imagePath string
}

type treatmentDay struct {
	userName  string
	petNumber int
	treatDay  string
}

func initDB() *sql.DB {
	psqlConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConnection)
	if err != nil {
		log.Panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	return db
}

func createTables() {
	db := initDB()
	defer db.Close()
	db.Query("CREATE TABLE IF NOT EXISTS pets (" +
		"username varchar(32) NOT NULL, " +
		"pet_number integer not null, " +
		"pet_name varchar(50) not null, " +
		"birthday char(10), " +
		"img_path varchar(35)," +
		"PRIMARY KEY(username,pet_number)" +
		")")
	db.Query("CREATE TABLE IF NOT EXISTS treatment_calendar (" +
		"username varchar(32) NOT NULL, " +
		"pet_number integer not null, " +
		"date_of_treatment char(10), " +
		"PRIMARY KEY(username,pet_number)" +
		")")
}

func addNewTreatmentDay(day treatmentDay) {
	db := initDB()
	defer db.Close()

	quertyInsertTreatmentDay := `INSERT INTO treatment_calendar 
    (username, pet_number, date_of_treatment) VALUES ($1, $2, $3)`

	_, err := db.Query(quertyInsertTreatmentDay, day.userName, day.petNumber, day.treatDay)
	if err != nil {
		log.Panic(err)
	}
}

func addNewPet(p pet) {
	db := initDB()
	defer db.Close()
	quertyInsertPet := `INSERT INTO pets 
    (username, pet_number, pet_name, birthday, img_path) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Query(quertyInsertPet, p.userName, p.petNumber, p.petName, p.birthday, p.imagePath)
	if err != nil {
		log.Panic(err)
	}
}
