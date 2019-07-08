package main

import (
	"fmt"
	"testing"

	_ "github.com/lib/pq"
)

func TestDB(t *testing.T) {
	dbinfo := fmt.Sprintf("host=%s port=5432 user=%s dbname=%s password=%s sslmode=disable",
		DB_HOST, DB_USER, DB_NAME, DB_PASSWORD)
	db, err := OpenDB(dbinfo)
	defer db.Close()
	if err != nil {
		t.Errorf("There was no database to connect to %v", err)
	}
}

