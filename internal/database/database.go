// Pacage database provides SQLite connection pooling and schema initialization.
package database

import (
	"database/sql"
	"fmt"
	_"log"
	_"embed" 
	_"modernc.org/sqlite"
)
// This looks like a comment, but there's no  space after // so the compiler treats it as the command itself.
//go:embed schema.sql
// When we run go build or go run, Go compiler opens schema.sql and assigns the text directly to the variable underneath which is schemaSQL
var schemaSQL string

// Open establishes a connection to the SQLite database at the specified path, dbPath in this case.
// verifies the connection via ping, and executes the intial schema DDL.
func Open(dbPath string) (*sql.DB, error){
	fmt.Println("This is a database.")

	db, err := sql.Open("sqlite",dbPath) // This is use to open a connection to the database pool
	if err != nil{
		return nil, err
	}
	//defer db.Close() 
	// This will close the connection  after the function returns
	// Doing defer in the function that opens it causes the function to close the database before returning, bad thing to do.
	err = db.Ping()
	if err != nil{
		return nil, err
	}
// Execute the embedded SQL, embeding is done during compile time not runtime
	if _, err := db.Exec(schemaSQL); err != nil{
		return nil, err
	}
	return db, err
}
