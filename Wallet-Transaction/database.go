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

func doTransaction(db *sql.DB, from string, fromName string, to string, toName string, amount uint64) bool {

	if(from==to){
		return false
	}
	
	row, err := db.Query("select * from amount where id=$1", from)

	if err != nil {
		log.Println(err)
	}

	var id string
    var balance uint64

	for row.Next() {
		row.Scan(&id, &balance)
	}

	if(balance<amount){
		return false
	}

	tx, err := db.Begin()
    if err != nil {
        return false
    }
	_, errFrom := tx.Exec("update amount set balance = balance - $1 where id = $2", amount, from)

	if errFrom != nil {
		tx.Rollback()
		log.Println(errFrom)
		return false
	}

	_, errTo := tx.Exec("update amount set balance = balance + $1 where id = $2", amount, to)

	if errTo != nil {
		tx.Rollback()
		return false
	}

	dt := time.Now()

	_, errTrans :=
		tx.Exec("insert into transactions(transactionTime,fromID,toID,toName,amount,fromName,isGenerated) values($1,$2,$3,$4,$5,$6,$7)",
			dt.Format("01-02-2006 15:04:05"), from, to, toName, amount, fromName,false)

	if errTrans != nil {
		tx.Rollback()
		return false
	}

	tx.Commit()

	return true
}

func addMoney(db *sql.DB, from string, fromName string, to string, toName string, amount uint64) bool {
	tx, err := db.Begin()
    if err != nil {
        return false
    }
	_, errTo := tx.Exec("update amount set balance = balance + $1 where id = $2", amount, to)
	if errTo != nil {
		tx.Rollback()
		return false
	}
	dt := time.Now()
	_, errTrans :=
		tx.Exec("insert into transactions(transactionTime,fromID,toID,toName,amount,fromName,isGenerated) values($1,$2,$3,$4,$5,$6,$7)",
			dt.Format("01-02-2006 15:04:05"), from, to, toName, amount, fromName,true)

	if errTrans != nil {
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

// MyState ...
type MyState struct {
	Balance      int
	Transactions []Transaction
}

// Transaction ...
type Transaction struct {
	From            interface{}
	To              interface{}
	TransactionID   interface{}
	TransactionTime interface{}
	ToName          interface{}
	Amount          interface{}
	FromName        interface{}
	isGenerated  	interface{}
}

func getMyState(db *sql.DB, number string) MyState {

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

	myState := MyState{state[number], getTransactions(db, number)}

	return myState
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

func getTransactions(sb *sql.DB, number string) []Transaction {

	log.Println("getting data of ", number)

	transactions := []Transaction{}

	row, err := db.Query("select * from transactions where fromid = $1 or toid = $1", number)

	if err != nil {
		log.Println(err)
	}

	for row.Next() {
		var transaction Transaction
		row.Scan(&transaction.TransactionID, &transaction.TransactionTime, &transaction.From, &transaction.To, &transaction.ToName, &transaction.FromName, &transaction.Amount ,&transaction.isGenerated)
		transactions = append(transactions, transaction)
    	log.Println(transaction.TransactionID, transaction.TransactionTime, transaction.From, transaction.To, transaction.ToName, transaction.FromName, transaction.Amount, transaction.isGenerated)

	}

	log.Println(transactions)

	return transactions
}

func getTransactionsBetweenObjects(sb *sql.DB, number1 string, number2 string) []Transaction {

	transactions := []Transaction{}

	row, err := db.Query("select * from transactions where (fromid = $1 or fromid = $2) and (toid = $1 or toid = $2)", number1, number2)

	if err != nil {
		log.Println(err)
	}

	for row.Next() {
		var transaction Transaction
		row.Scan(&transaction.TransactionID, &transaction.TransactionTime, &transaction.From, &transaction.To, &transaction.ToName, &transaction.FromName, &transaction.Amount, &transaction.isGenerated)
		transactions = append(transactions, transaction)

	}


	return transactions
}
