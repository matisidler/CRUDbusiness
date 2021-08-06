package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/matisidler/CRUDbusiness/pkg/clients"
)

const (
	MigrateClient = `CREATE TABLE IF NOT EXISTS clients(
		id SERIAL NOT NULL,
		name VARCHAR(60) NOT NULL,
		country VARCHAR(50) NOT NULL,
		postal_code VARCHAR(20) NOT NULL,
		comment VARCHAR(150),
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP,
		CONSTRAINT clients_id_pk PRIMARY KEY (id))`
	InsertClient  = `INSERT INTO clients(name,country,postal_code,comment) VALUES($1,$2,$3,$4) RETURNING id`
	getAllClients = `SELECT * FROM clients`
	getByIdClient = `SELECT * FROM clients WHERE id = $1`
	updateClient  = `UPDATE clients SET name = $1, country = $2, postal_code= $3, comment = $4, updated_at = now() WHERE id = $5`
	deleteClient  = `DELETE FROM clients WHERE id = $1`
)

type PsqlClient struct {
	db *sql.DB
}

func NewPsqlClient(db *sql.DB) *PsqlClient {
	return &PsqlClient{db}
}

func (p *PsqlClient) Migrate() error {
	db, err := p.db.Prepare(MigrateClient)
	if err != nil {
		return err
	}
	_, err = db.Exec()
	if err != nil {
		return err
	}
	fmt.Println("succesfuly client migration.")
	defer db.Close()
	return nil
}

func (p *PsqlClient) Insert(m *clients.Model) error {
	db, err := p.db.Prepare(InsertClient)
	if err != nil {
		return err
	}
	defer db.Close()
	m.CreatedAt = time.Now()
	err = db.QueryRow(stringToNull(m.Name), stringToNull(m.Country), stringToNull(m.PostalCode), stringToNull(m.Comment)).Scan(&m.ID)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", m)
	fmt.Println("inserted client")
	return nil
}

func (p *PsqlClient) GetAll() ([]*clients.Model, error) {
	db, err := p.db.Prepare(getAllClients)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var models []*clients.Model
	res, err := db.Query()
	if err != nil {
		return nil, err
	}
	for res.Next() {
		m := clients.Model{}
		res.Scan(&m.ID, &m.Name, &m.Country, &m.PostalCode, &strNull, &m.CreatedAt, &timeNull)
		m.UpdatedAt = timeNull.Time
		m.Comment = strNull.String
		models = append(models, &m)
	}
	return models, nil
}

func (p *PsqlClient) GetById(id uint) (*clients.Model, error) {
	stmt, err := p.db.Prepare(getByIdClient)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	m := &clients.Model{}
	err = stmt.QueryRow(id).Scan(&m.ID, &m.Name, &m.Country, &m.PostalCode, &strNull, &m.CreatedAt, &timeNull)
	m.UpdatedAt = timeNull.Time
	m.Comment = strNull.String
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (p *PsqlClient) Update(m *clients.Model) error {
	db, err := p.db.Prepare(updateClient)
	if err != nil {
		return err
	}
	defer db.Close()
	res, err := db.Exec(&m.Name, &m.Country, &m.PostalCode, stringToNull(m.Comment), &m.ID)
	if err != nil {
		return err
	}
	if res, _ := res.RowsAffected(); res != 1 {
		return errors.New("more than 1 (or 0) rows modified")
	}
	fmt.Println("client updated succesfuly")
	return nil
}

func (p *PsqlClient) Delete(id uint) error {
	stmt, err := p.db.Prepare(deleteClient)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows != 1 {
		return errors.New("error: more than 1 (or 0) rows modified")
	}
	fmt.Printf("client with id = %d deleted succesfuly\n", id)
	return nil
}
