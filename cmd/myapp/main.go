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

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	database *databaseConn
}

type userInfoStruct struct {
	Email  string
	UserId int
}

var userInfo userInfoStruct

func main() {

	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "root:Admin123###@tcp(localhost:3308)/dekkoto?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		database: &databaseConn{DB: db},
	}

	// Initialize Gin router
	router := gin.Default()
	// load all the html files
	router.LoadHTMLGlob("ui/html/*")
	// load all the static files
	router.Static("/static", "./ui/static")

	// load all the images and videos from ./userUploadDatas
	router.StaticFS("/userUploadDatas", http.Dir("./userUploadDatas"))

	// call the routes function from routes.go
	app.routes(router)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  router,
	}
	infoLog.Printf("Starting server on localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
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
