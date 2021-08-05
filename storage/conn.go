package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

type Driver string

const (
	Postgres Driver = "Postgres"
)

func NewConnection(d Driver) {
	switch d {
	case Postgres:
		psqlConn()
	default:
		log.Fatal("can't found the db driver")
	}
}

func psqlConn() *sql.DB {
	once.Do(func() {
		db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/business?sslmode=disable")
		if err != nil {
			log.Fatalf("can't open DB: %v\n", err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatalf("can't ping to DB: %v\n", err)
		}
		fmt.Println("connected to postgres")
	})
	return db
}
