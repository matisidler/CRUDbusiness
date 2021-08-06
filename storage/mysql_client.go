package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/matisidler/CRUDbusiness/pkg/clients"
)

const (
	myMigrateClient = `CREATE TABLE IF NOT EXISTS clients(
		id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
		name VARCHAR(60) NOT NULL,
		country VARCHAR(50) NOT NULL,
		postal_code VARCHAR(20) NOT NULL,
		comment VARCHAR(150),
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP)`
	myInsertClient  = `INSERT INTO clients(name,country,postal_code,comment) VALUES(?,?,?,?)`
	mygetAllClients = `SELECT * FROM clients`
	mygetByIdClient = `SELECT * FROM clients WHERE id = ?`
	myupdateClient  = `UPDATE clients SET name = ?, country = ?, postal_code= ?, comment = ?, updated_at = now() WHERE id = ?`
	mydeleteClient  = `DELETE FROM clients WHERE id = ?`
)

type MySqlClient struct {
	db *sql.DB
}

func NewMySqlClient(db *sql.DB) *MySqlClient {
	return &MySqlClient{db}
}

func (p *MySqlClient) Migrate() error {
	db, err := p.db.Prepare(myMigrateClient)
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

func (p *MySqlClient) Insert(m *clients.Model) error {
	db, err := p.db.Prepare(myInsertClient)
	if err != nil {
		return err
	}
	defer db.Close()
	m.CreatedAt = time.Now()
	res, err := db.Exec(stringToNull(m.Name), stringToNull(m.Country), stringToNull(m.PostalCode), stringToNull(m.Comment))
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	m.ID = uint(id)
	fmt.Printf("%+v\n", m)
	fmt.Println("inserted client")
	return nil
}

func (p *MySqlClient) GetAll() ([]*clients.Model, error) {
	db, err := p.db.Prepare(mygetAllClients)
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

func (p *MySqlClient) GetById(id uint) (*clients.Model, error) {
	stmt, err := p.db.Prepare(mygetByIdClient)
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

func (p *MySqlClient) Update(m *clients.Model) error {
	db, err := p.db.Prepare(myupdateClient)
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

func (p *MySqlClient) Delete(id uint) error {
	stmt, err := p.db.Prepare(mydeleteClient)
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
	fmt.Printf("sale with id = %d deleted succesfuly\n", id)
	return nil
}
