package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/matisidler/CRUDbusiness/pkg/products"
)

const (
	MigrateProduct = `CREATE TABLE IF NOT EXISTS products(
		id SERIAL NOT NULL,
		name VARCHAR(60) NOT NULL,
		price INT NOT NULL,
		observations VARCHAR(100),
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP,
		CONSTRAINT products_id_pk PRIMARY KEY (id))`
	InsertProduct  = `INSERT INTO products(name, price, observations) VALUES($1,$2,$3) RETURNING id`
	getAllProducts = `SELECT * FROM products`
	getByIdProduct = getAllProducts + ` WHERE id = $1`
	updateProduct  = `UPDATE products SET name = $1, price = $2, observations = $3, updated_at = now() WHERE id = $4`
	deleteProduct  = `DELETE FROM products WHERE id = $1`
)

type PsqlProduct struct {
	db *sql.DB
}

func NewPsqlProduct(db *sql.DB) *PsqlProduct {
	return &PsqlProduct{db}
}

var timeNull = sql.NullTime{}
var strNull = sql.NullString{}

func (p *PsqlProduct) Migrate() error {
	db, err := p.db.Prepare(MigrateProduct)
	if err != nil {
		return err
	}
	_, err = db.Exec()
	if err != nil {
		return err
	}
	fmt.Println("succesfuly product migration.")
	defer db.Close()
	return nil
}

func (p *PsqlProduct) Insert(m *products.Model) error {
	db, err := p.db.Prepare(InsertProduct)
	if err != nil {
		return err
	}
	defer db.Close()
	m.CreatedAt = time.Now()
	err = db.QueryRow(stringToNull(m.Name), intToNull(int(m.Price)), stringToNull(m.Observations)).Scan(&m.ID)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", m)
	fmt.Println("inserted product")
	return nil
}

func (p *PsqlProduct) GetAll() ([]*products.Model, error) {
	stmt, err := p.db.Prepare(getAllProducts)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	var models []*products.Model
	for res.Next() {
		m := &products.Model{}
		err := res.Scan(&m.ID, &m.Name, &m.Price, &strNull, &m.CreatedAt, &timeNull)
		m.UpdatedAt = timeNull.Time
		m.Observations = strNull.String
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, nil
}

func (p *PsqlProduct) GetById(id int64) (*products.Model, error) {
	stmt, err := p.db.Prepare(getByIdProduct)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	m := &products.Model{}
	err = stmt.QueryRow(id).Scan(&m.ID, &m.Name, &m.Price, &strNull, &m.CreatedAt, &timeNull)
	m.UpdatedAt = timeNull.Time
	m.Observations = strNull.String
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (p *PsqlProduct) Update(m *products.Model) error {
	stmt, err := p.db.Prepare(updateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(stringToNull(m.Name), intToNull(int(m.Price)), stringToNull(m.Observations), intToNull(int(m.ID)))
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return errors.New("more than 1 (or 0) rows modified")
	}
	fmt.Println("product updated.")

	return nil
}

func (p *PsqlProduct) Delete(id uint) error {
	stmt, err := p.db.Prepare(deleteProduct)
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
	fmt.Printf("product with id = %d deleted succesfuly\n", id)
	return nil
}
