package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Connect ...
func Connect() *sql.DB {
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

func doTransaction(db *sql.DB, from string, to string, amount uint64) bool {
	_, errFrom := db.Exec("update amount set balance = balance - $1 where id = $2", amount, from)

	if errFrom != nil {
		log.Println(errFrom)
		return false
	}

	_, errTo := db.Exec("update amount set balance = balance + $1 where id = $2", amount, to)

	if errTo != nil {
		db.Exec("update amount set balance = balance + $1 where id = $2", amount, from)
		return false
	}

	dt := time.Now()

	_, errTrans := db.Exec("insert into transactions(transactionTime,fromID,toID,amount) values($1,$2,$3,$4)", dt.Format("01-02-2006 15:04:05"), from, to, amount)

	if errTrans != nil {
		db.Exec("update amount set balance = balance + $1 where id = $2", amount, from)
		db.Exec("update amount set balance = balance - $1 where id = $2", amount, to)
		return false
	}

	return true
}

func getMyState(db *sql.DB, number string) map[string]int {

	state := map[string]int{}

	row, err := db.Query("select * from amount where id=$1", number)

	if err != nil {
		log.Println(err)
	}

	for row.Next() {
		var id string
		var balance int
		row.Scan(&id, &balance)
		state[id] = balance
	}

	return state
}

func getState(db *sql.DB) map[string]int {

	state := map[string]int{}

	row, err := db.Query("select * from amount")

	if err != nil {
		log.Println(err)
	}

	for row.Next() {
		var id string
		var balance int
		row.Scan(&id, &balance)
		state[id] = balance

	}

	return state

}

// Transaction ...
type Transaction struct {
	FromID          interface{}
	ToID            interface{}
	TransactionID   interface{}
	TransactionTime interface{}
	Amount          interface{}
}

func getTransaction(sb *sql.DB, number string) []Transaction {

	transactions := []Transaction{}

	row, err := db.Query("select * from transactions")

	if err != nil {
		log.Println(err)
	}

	for row.Next() {
		log.Println("Row 1")
		var transaction Transaction
		row.Scan(&transaction.TransactionID, &transaction.TransactionTime, &transaction.FromID, &transaction.ToID, &transaction.Amount)
		transactions = append(transactions, transaction)
	}

	return transactions
}
