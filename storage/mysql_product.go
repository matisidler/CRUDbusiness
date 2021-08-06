package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/matisidler/CRUDbusiness/pkg/products"
)

const (
	myMigrateProduct = `CREATE TABLE IF NOT EXISTS products(
		id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
		name VARCHAR(60) NOT NULL,
		price INT NOT NULL,
		observations VARCHAR(100),
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP
		)`
	myInsertProduct  = `INSERT INTO products(name, price, observations) VALUES(?,?,?)`
	mygetAllProducts = `SELECT * FROM products`
	mygetByIdProduct = getAllProducts + ` WHERE id = ?`
	myupdateProduct  = `UPDATE products SET name = ?, price = ?, observations = ?, updated_at = now() WHERE id = ?`
	mydeleteProduct  = `DELETE FROM products WHERE id = ?`
)

type MySqlProduct struct {
	db *sql.DB
}

func NewMySqlProduct(db *sql.DB) *MySqlProduct {
	return &MySqlProduct{db}
}

func (p *MySqlProduct) Migrate() error {
	db, err := p.db.Prepare(myMigrateProduct)
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

func (p *MySqlProduct) Insert(m *products.Model) error {
	db, err := p.db.Prepare(myInsertProduct)
	if err != nil {
		return err
	}
	defer db.Close()
	m.CreatedAt = time.Now()
	res, err := db.Exec(stringToNull(m.Name), intToNull(int(m.Price)), stringToNull(m.Observations))
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	m.ID = uint(id)
	fmt.Printf("%+v\n", m)
	fmt.Println("inserted product")
	return nil
}

func (p *MySqlProduct) GetAll() ([]*products.Model, error) {
	stmt, err := p.db.Prepare(mygetAllProducts)
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

func (p *MySqlProduct) GetById(id int64) (*products.Model, error) {
	stmt, err := p.db.Prepare(mygetByIdProduct)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	m := &products.Model{}
	err = stmt.QueryRow(id).Scan(&m.ID, &m.Name, &m.Price, &strNull, &m.CreatedAt, &timeNull)
	if err != nil {
		return nil, err
	}
	m.UpdatedAt = timeNull.Time
	m.Observations = strNull.String
	return m, nil
}

func (p *MySqlProduct) Update(m *products.Model) error {
	stmt, err := p.db.Prepare(myupdateProduct)
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

func (p *MySqlProduct) Delete(id uint) error {
	stmt, err := p.db.Prepare(mydeleteProduct)
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
