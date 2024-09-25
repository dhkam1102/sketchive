package main

import (
	"database/sql"
	"fmt"
	"log"
	"sketchive/internal/db"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//dsn: Data Source Name
	dsn := "root:@tcp(127.0.0.1:3306)/sketchive"
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Could not grad the connvection at first place", err)
	}

	// Ping() checks if the connection is alive
	err = database.Ping()
	if err != nil {
		log.Fatal("Lost Database connection failed:", err)
	} else {
		fmt.Println("Successfully connected to database!")
	}

	db.SetDB(database)
}
