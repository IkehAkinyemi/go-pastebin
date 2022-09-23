package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// application struct hold the application-wide dependencies
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// configuration
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "dev:Akinyemi1234@@/snippetbox?parseTime=true", "MySQL data source name")
	
	flag.Parse()

	// Create INFO and ERROR loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Pass DSN to openDB to keep main succinct
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Initialize server struct to support custom error logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Use the http.ListenAndServe() function to start a new web server.
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
