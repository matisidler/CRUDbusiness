package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/matisidler/CRUDbusiness/pkg/sales"
)

const (
	MigrateSale = `CREATE TABLE IF NOT EXISTS sales(
		id SERIAL NOT NULL,
		id_client INT NOT NULL,
		id_product INT NOT NULL,
		product_price INT,
		quantity INT NOT NULL,
		income INT,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP,
		CONSTRAINT sales_id_pk PRIMARY KEY (id),
		CONSTRAINT sales_id_client_fk FOREIGN KEY (id_client) REFERENCES clients (id),
		CONSTRAINT sales_id_product_fk FOREIGN KEY (id_product) REFERENCES products (id))`
	InsertSale  = `INSERT INTO sales(id_client,id_product,quantity) VALUES($1,$2,$3) RETURNING id`
	InsertSale2 = `UPDATE sales SET product_price = products.Price FROM products WHERE products.id = sales.id_product AND sales.id = $1`
	InsertSale3 = `UPDATE sales SET income = product_price * quantity`
	getAllSales = `SELECT * FROM sales`
	getByIdSale = getAllSales + ` WHERE id = $1`
	updateSale  = `UPDATE sales SET id_client = $1, id_product = $2, quantity = $3, updated_at = now() WHERE id = $4`
	deleteSale  = `DELETE FROM sales WHERE id = $1`
)

type PsqlSale struct {
	db *sql.DB
}

func NewPsqlSale(db *sql.DB) *PsqlSale {
	return &PsqlSale{db}
}

func (p *PsqlSale) Migrate() error {
	db, err := p.db.Prepare(MigrateSale)
	if err != nil {
		return err
	}
	_, err = db.Exec()
	if err != nil {
		return err
	}
	fmt.Println("succesfuly sales migration.")
	defer db.Close()
	return nil
}

func (p *PsqlSale) Insert(m *sales.Model) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	err = funcInsertSale1(tx, m)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("InsertSale1: %v\n", err)
	}
	err = funcInsertSale2(tx, m.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("InsertSale2: %v\n", err)
	}
	err = funcInsertSale3(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("InsertSale3 %v\n", err)
	}
	tx.Commit()
	fmt.Println("Insert completed")
	return nil

}

func funcInsertSale1(tx *sql.Tx, m *sales.Model) error {
	stmt, err := tx.Prepare(InsertSale)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(intToNull(int(m.IdClient)), intToNull(int(m.IdProduct)), intToNull(int(m.Quantity))).Scan(&m.ID)
	if err != nil {
		return err
	}
	fmt.Println("insert number 1 completed")
	return nil
}

func funcInsertSale2(tx *sql.Tx, i uint) error {
	stmt, err := tx.Prepare(InsertSale2)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(i)
	if err != nil {
		return err
	}
	fmt.Println("insert number 2 completed")
	return nil
}

func funcInsertSale3(tx *sql.Tx) error {
	stmt, err := tx.Prepare(InsertSale3)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	fmt.Println("insert number 3 completed")
	return nil
}

func (p *PsqlSale) GetAll() ([]*sales.Model, error) {
	stmt, err := p.db.Prepare(getAllSales)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Query()
	var models []*sales.Model
	for res.Next() {
		m := &sales.Model{}
		err := res.Scan(&m.ID, &m.IdClient, &m.IdProduct, &m.ProductPrice, &m.Quantity, &m.Income, &m.CreatedAt, &timeNull)
		m.UpdatedAt = timeNull.Time
		if err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, nil
}

func (p *PsqlSale) GetById(id uint) (*sales.Model, error) {
	stmt, err := p.db.Prepare(getByIdSale)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	m := &sales.Model{}
	err = stmt.QueryRow(id).Scan(&m.ID, &m.IdClient, &m.IdProduct, &m.ProductPrice, &m.Quantity, &m.Income, &m.CreatedAt, &timeNull)
	m.UpdatedAt = timeNull.Time
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (p *PsqlSale) Update(m *sales.Model) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(updateSale)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(&m.IdClient, &m.IdProduct, &m.Quantity, &m.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		tx.Rollback()
		return err
	}
	err = funcInsertSale2(tx, m.ID)
	if err != nil {
		return err
	}
	err = funcInsertSale3(tx)
	if err != nil {
		return err
	}
	tx.Commit()
	fmt.Println("sale updated succesfuly")
	return nil

}
func (p *PsqlSale) Delete(id uint) error {
	stmt, err := p.db.Prepare(deleteSale)
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
