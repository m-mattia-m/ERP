package db

import (
	"database/sql"
	"fmt"
)

const (
	host     string = "localhost"
	port     int    = 6603
	username string = "systemuser"
	password string = "password"
	database string = "ERP"
)

func InitDB() {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)
	// db, err := sql.Open("mysql", "root:6HmFbaZH2X3Z0jXjruqD@tcp(localhost:6603)/ERP")

	if err != nil {
		// Error: cant connect to DB
		fmt.Println(err.Error())
	} else if err = db.Ping(); err != nil {
		// Error: lost connect to DB
		fmt.Println(err.Error())
	}
	defer db.Close()
	var createUsersQuers string = `CREATE TABLE IF NOT EXISTS users (` +
		`RecordId int AUTO_INCREMENT PRIMARY KEY NOT NULL,` +
		`Id varchar(255) NOT NULL,` +
		`Firstname varchar(255) NOT NULL,` +
		`Lastname varchar(255) NOT NULL,` +
		`Username varchar(255) NOT NULL,` +
		`Email varchar(255) NOT NULL,` +
		`Password varchar(255) NOT NULL,` +
		`Role varchar(255) NOT NULL` +
		`);`
	res, err := db.Exec(createUsersQuers)
	if err != nil {
		fmt.Println("[DB-Creation] Can't create users table \t--> " + err.Error())
	} else {
		fmt.Print("[DB-Creation] Create users table was successfully \t--> ")
		fmt.Print(res)
		fmt.Print("\n")
	}

	var createCustomerQuers string = `CREATE TABLE IF NOT EXISTS customer (` +
		`RecordId int AUTO_INCREMENT PRIMARY KEY NOT NULL,` +
		`Id varchar(255) NOT NULL,` +
		`Firstname varchar(255) NOT NULL,` +
		`Lastname varchar(255) NOT NULL,` +
		`Street varchar(255) NOT NULL,` +
		`StreetNr int(10) NOT NULL,` +
		`Postcode int(4) NOT NULL,` +
		`City varchar(255) NOT NULL,` +
		`Email varchar(255) NOT NULL,` +
		`Telefon int(16) NOT NULL` +
		`);`
	res, err = db.Exec(createCustomerQuers)
	if err != nil {
		fmt.Println("[DB-Creation] Can't create customers table \t--> " + err.Error())
	} else {
		fmt.Print("[DB-Creation] Create customers table was successfully \t--> ")
		fmt.Print(res)
		fmt.Print("\n")
	}

	var createReportsQuers string = `CREATE TABLE IF NOT EXISTS reports (` +
		`RecordId int AUTO_INCREMENT PRIMARY KEY NOT NULL,` +
		`Id varchar(255) NOT NULL,` +
		`Title varchar(255) NOT NULL,` +
		`Description text NOT NULL,` +
		`Datum datetime NOT NULL,` +
		`Dauer double NOT NULL,` +
		`KundenId varchar(255) NOT NULL,` +
		`CreaterId varchar(255) NOT NULL` +
		`);`
	res, err = db.Exec(createReportsQuers)
	if err != nil {
		fmt.Println("[DB-Creation] Can't create reports table \t--> " + err.Error())
	} else {
		fmt.Print("[DB-Creation] Create reports table was successfully \t--> ")
		fmt.Print(res)
		fmt.Print("\n")
	}
	defer db.Close()
}

func RunSqlQueryWithReturn(query string) (*sql.Rows, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		// Error: cant connect to DB
		return nil, err
	} else if err = db.Ping(); err != nil {
		// Error: lost connect to DB
		return nil, err
	}
	defer db.Close()
	res, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer db.Close()
	return res, nil
}

func RunSqlQueryWithSingeReturn(query string) (*sql.Rows, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		// Error: cant connect to DB
		return nil, err
	} else if err = db.Ping(); err != nil {
		// Error: lost connect to DB
		return nil, err
	}
	defer db.Close()
	res, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer db.Close()
	return res, nil
}

func RunSqlQueryWithoutReturn(query string) (bool, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	db, err := sql.Open("mysql", connString)

	if err != nil {
		// Error: cant connect to DB
		return false, err
	} else if err = db.Ping(); err != nil {
		// Error: lost connect to DB
		return false, err
	}
	defer db.Close()
	res, err := db.Exec(query)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	fmt.Println(res)
	defer db.Close()
	return true, nil
}
