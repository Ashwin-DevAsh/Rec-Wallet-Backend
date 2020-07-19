package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func connect() *sql.DB {
	const (
		host     = "status-database"
		port     = 5432
		user     = "postgres"
		password = "2017PASS"
		dbname   = "Rec_Wallet"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Println(err)
	}

	return db
}

func updateOnline(db *sql.DB, id string, socketID string, fcmToken string, isOnline bool) {
	insertStatement := `update info set isonline=$4 , socketid=$2 , fcmToken = $3  where id=$1`
	_, err := db.Exec(insertStatement, id, socketID, fcmToken, isOnline)
	log.Println("UpdatedOnline", id,socketID)
}

func updateOffline(db *sql.DB, socketID string) {
	insertStatement := `update info set socketid=null , isonline=false where socketid=$1`
	_, err := db.Exec(insertStatement, socketID)
	log.Println("UpdatedOffline", socketID)
}
