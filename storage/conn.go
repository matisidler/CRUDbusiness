package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/matisidler/CRUDbusiness/pkg/clients"
	"github.com/matisidler/CRUDbusiness/pkg/products"
	"github.com/matisidler/CRUDbusiness/pkg/sales"
)

var (
	db   *sql.DB
	once sync.Once
)

type Driver string

const (
	Postgres Driver = "Postgres"
	MySQL    Driver = "MySQL"
)

func NewConnection(d Driver) {
	switch d {
	case Postgres:
		newPqDB()
	case MySQL:
		newMySqlDB()
	default:
		log.Fatal("can't found the db driver")
	}
}

func newPqDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("postgres", "postgres://postgres:password@localhost:5432/business?sslmode=disable")
		if err != nil {
			log.Fatalf("can't open DB %v", err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatalf("can't do ping: %v", err)
		}
		fmt.Println("connected to postgres.")
	})
	return db
}
func newMySqlDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/business?parseTime=true")
		if err != nil {
			log.Fatalf("can't open DB %v", err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatalf("can't do ping: %v", err)
		}
		fmt.Println("connected to mysql.")
	})
	return db
}

func DAOClient(d Driver) (clients.Storage, error) {
	switch d {
	case Postgres:
		return NewPsqlClient(db), nil
	case MySQL:
		return NewMySqlClient(db), nil
	default:
		return nil, errors.New("not allowed driver.")
	}
}

func DAOProduct(d Driver) (products.Storage, error) {
	switch d {
	case Postgres:
		return NewPsqlProduct(db), nil
	case MySQL:
		return NewMySqlProduct(db), nil
	default:
		return nil, errors.New("not allowed driver.")
	}
}

func DAOSale(d Driver) (sales.Storage, error) {
	switch d {
	case Postgres:
		return NewPsqlSale(db), nil
	case MySQL:
		return NewMySqlSale(db), nil
	default:
		return nil, errors.New("not allowed driver.")
	}
}

func stringToNull(m string) *sql.NullString {
	strNull := sql.NullString{}
	if m == "" {
		strNull.Valid = false

	} else {
		strNull.Valid = true
		strNull.String = m
	}
	return &strNull
}

func intToNull(i int) *sql.NullInt64 {
	intNull := sql.NullInt64{}
	if i == 0 {
		intNull.Valid = false
	} else {
		intNull.Valid = true
		intNull.Int64 = int64(i)
	}
	return &intNull
}
