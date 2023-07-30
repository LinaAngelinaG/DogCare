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
	chatId    int64
	petNumber int16
	petName   string
	birthday  string
	imagePath string
}

type treatmentDay struct {
	chatId    string
	petNumber int
	treatDay  string
}

func initDB() *sql.DB {
	psqlConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s "+
		"dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlConnection)
	checkError(err)
	err = db.Ping()
	checkError(err)
	return db
}

func createTables() {
	db := initDB()
	defer db.Close()
	db.Query("CREATE TABLE IF NOT EXISTS pets (" +
		"chatId integer NOT NULL, " +
		"pet_number integer not null, " +
		"pet_name varchar(50) not null, " +
		"birthday char(10), " +
		"img_path varchar(35)," +
		"PRIMARY KEY(username,pet_number)" +
		")")
	db.Query("CREATE TABLE IF NOT EXISTS treatment_calendar (" +
		"chatId integer NOT NULL, " +
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

	_, err := db.Query(quertyInsertTreatmentDay, day.chatId, day.petNumber, day.treatDay)
	checkError(err)
}

func getMaxPetNumber(charId int64) int16 {
	db := initDB()
	defer db.Close()
	quertyInsertPet := `SELECT pet_number FROM pets 
    ORDER BY pet_number DESC LIMIT 1`

	res := int16(1)

	rows, err := db.Query(quertyInsertPet)

	checkError(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&res)
		checkError(err)
	}
	return res
}

func addNewPet(p pet) {
	db := initDB()

	defer db.Close()
	quertyInsertPet := `INSERT INTO pets 
    (chatId, pet_number, pet_name, birthday, img_path) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Query(quertyInsertPet, p.chatId, p.petNumber, p.petName, p.birthday, p.imagePath)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
