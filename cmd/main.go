package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/obynonwane/my_blockchain_prototype/cmd/database"
	"github.com/obynonwane/my_blockchain_prototype/cmd/logger"
)

const publicPort = "8080"
const privatePort = "8081"
const webPort = "8082"

var counts int64

type Config struct {
	DB     *sql.DB
	Models database.Models
}

func main() {

	// Initialize logger
	logger.Init()

	//Connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to Postgres")
	}

	//setup config
	app := Config{
		DB:     conn,
		Models: database.New(conn),
	}

	// Start the service listening on public port.
	go func() {
		// define http server
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", publicPort),
			Handler: app.publicRoutes(),
		}

		// start the server
		err := srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}

	}()

	// Start the service listening on private port.
	go func() {
		// start second server port
		// define http server
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", privatePort),
			Handler: app.privateRoutes(),
		}

		// start the server
		err := srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}

	}()

	// Start the service listening on web port.
	go func() {
		// start second server port
		// define http server
		srv := &http.Server{
			Addr:    fmt.Sprintf(":%s", webPort),
			Handler: app.webRoutes(),
		}

		// start the server
		err := srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}

	}()

	// Prevent main from exiting
	select {}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {

	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")
	log.Println("db user", dbUser)

	// Construct the DSN string
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	log.Printf("%s", dsn)
	// dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres ...")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("backing off for 2 seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}
