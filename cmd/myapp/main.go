// Package main is the entry point of the application.
// It sets up the logging, database connection, and starts the HTTP server.
package main

import (
	"database/sql"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// application struct holds the application-wide dependencies, such as loggers and database connection.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	database *databaseConn
}

// userInfoStruct struct holds the user information.
type userInfoStruct struct {
	Email  string
	UserId int
}

// Global variable userInfo holds the user information.
var userInfo userInfoStruct

// The main function is the entry point of the application.
// It parses the command-line flags, sets up the logging, opens the database connection,
// initializes the application struct, sets up the routes, and starts the HTTP server.
func main() {

	// Parse the command-line flags.
	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "root:Admin123###@tcp(localhost:3308)/dekkoto?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Set up the loggers.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Open the database connection.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize the application struct.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		database: &databaseConn{DB: db},
	}

	// Initialize Gin router, load all the html files, static files and user uploaded data.
	router := gin.Default()
	router.LoadHTMLGlob("ui/html/*")
	router.Static("/static", "./ui/static")
	router.StaticFS("/userUploadDatas", http.Dir("./userUploadDatas"))

	// Set up the routes.
	app.routes(router)

	// Start the HTTP server.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  router,
	}
	infoLog.Printf("Starting server on localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// openDB function opens a new database connection.
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
