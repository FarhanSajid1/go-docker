package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// database attributes
var DB_USER = getEnv("PGUSER", "postgres")
var DB_PASSWORD = getEnv("PGPASSWORD", "postgres")
var DB_NAME = getEnv("PGDATABASE", "postgres")
var DB_HOST = getEnv("PGHOST", "localhost")

var db *sql.DB
var err error

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func OpenDB(dbinfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbinfo)
	return db, err
}

func DatabaseSetup() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable",
		DB_HOST, DB_USER, DB_NAME, DB_PASSWORD)

	db, err := OpenDB(dbinfo)
	checkErr(err)

	// close the database at the end
	defer db.Close()
	err1 := db.Ping()
	checkErr(err1)
	fmt.Println("Server connected")
	return db
}

func All(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Returning all the records from new table")
	rows, err := db.Query("Select name FROM nba;")
	defer rows.Close()

	var s, name string
	s = "Retrieved records:\n"
	for rows.Next() {
		err = rows.Scan(&name) // copies the data into the empty var name
		checkErr(err)
		s += name + "\n" // we are going to add to s with a new line
	}
	fmt.Fprintln(w, s)

}

func inserting(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`insert into nba values ('Lebron James'), ('Kobe Bryant');`)
	checkErr(err)
	res, err := stmt.Exec()
	checkErr(err)
	n, err := res.RowsAffected()
	checkErr(err)
	fmt.Fprintf(w, "Inserted records, affected %v rows", n)
}

var visited bool

func CreateTable(w http.ResponseWriter, r *http.Request) {
	if !visited {
		stmt, err := db.Prepare(`CREATE TABLE nba (name varchar(20));`)
		if err != nil {
			fmt.Fprintf(w, "%v", err)
		}
		// execute the statement
		res, err := stmt.Exec()
		checkErr(err)
		fmt.Fprintf(w, "Creating a new table %v", res)
		visited = true
	} else {
		http.Redirect(w, r, "/all/", 302)
	}
}

func deleteRecords(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`DELETE FROM nba WHERE name='Lebron James';`)
	checkErr(err)
	_, err1 := stmt.Exec()
	checkErr(err1)
	http.Redirect(w, r, "/all/", 302)
}

func updateRecords(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare(`UPDATE nba set name='farhan sajid'
	WHERE name='Kobe Bryant';`)
	checkErr(err)
	_, err1 := stmt.Exec()
	checkErr(err1)
	http.Redirect(w, r, "/all/", 302)
}
func main() {
	// Sprintf formats with the placeholders and returns as a string
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	db, err = sql.Open("postgres", dbinfo)
	checkErr(err)

	// close the database at the end
	defer db.Close()
	err1 := db.Ping()
	checkErr(err1)
	fmt.Println("Server connected")
	http.HandleFunc("/all/", All)
	http.HandleFunc("/create/", CreateTable)
	http.HandleFunc("/delete/", deleteRecords)
	http.HandleFunc("/insert/", inserting)
	http.HandleFunc("/update/", updateRecords)
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
	// function to handle the errors

	if err != nil {
		panic(err)
	}
}
