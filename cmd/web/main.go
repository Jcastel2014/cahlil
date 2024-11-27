package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
	"sysadmin.com/final/pkg/models/postgresql"
)

func ConnectToDatabase(settings serverConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	fmt.Println("Database Connection established")
	return db, nil
}

type serverConfig struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	user     *postgresql.UserModel
	errorLog *log.Logger
	infoLog  *log.Logger
	session  *sessions.Session
	config   serverConfig
}

func main() {
	var settings serverConfig

	flag.IntVar(&settings.port, "port", 4000, "Server Port")
	flag.StringVar(&settings.env, "env", "development", "Environment(Development|Staging|Production)")
	flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://users:password@db/users?sslmode=disable", "PostgreSQL DSN")
	secret := flag.String("secret", "p7Mhd+qQamgHsS*+8Tg7mNXtcjvu@egz", "Secret Key")
	flag.Parse()
	log.Printf("Starting server on port: %d", settings.port)
	log.Printf("Environment: %s", settings.env)
	log.Printf("Database DSN: %s", settings.db.dsn)
	log.Printf("Secret Key: %s", *secret)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := ConnectToDatabase(settings)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	session := sessions.New([]byte(*secret))
	session.Lifetime = 1 * time.Minute
	session.Secure = true

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		session:  session,
		user: &postgresql.UserModel{
			DB: db,
		},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", settings.port),
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on port %s", srv.Addr)
	err = srv.ListenAndServe()
	srv.ErrorLog.Fatal(err)
}
