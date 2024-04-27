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

	_ "github.com/lib/pq"
	"sulfurAuth.net/internal/mailer"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	logger *log.Logger
	DB     *sql.DB
	mailer mailer.Mailer
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:sadboys133722813@localhost:5432/reg?sslmode=disable", "PostgreSQL DSN")

	//smtp
	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smt-name", "209a14c24fb89e", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "6da2405eaedd26", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Temirlan <no-reply@sulfurAuth.net>", "SMT sender")

	flag.Parse()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// db := db.GetConnetWithMongo()
	DB, err := openDB(cfg)
	if err != nil {
		log.Fatal("Error connecting db")
	}
	defer DB.Close()
	app := &application{
		config: cfg,
		logger: logger,
		DB:     DB,
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/registration", app.CreateUser)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
