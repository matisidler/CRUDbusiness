package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/matisidler/CRUDbusiness/pkg/sales"
)

const (
	myMigrateSale = `CREATE TABLE IF NOT EXISTS sales(
		id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
		id_client INT NOT NULL,
		id_product INT NOT NULL,
		product_price INT,
		quantity INT NOT NULL,
		income INT,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP,
		CONSTRAINT sales_id_client_fk FOREIGN KEY (id_client) REFERENCES clients (id),
		CONSTRAINT sales_id_product_fk FOREIGN KEY (id_product) REFERENCES products (id))`
	myInsertSale  = `INSERT INTO sales(id_client,id_product,quantity) VALUES(?,?,?)`
	myInsertSale2 = `UPDATE sales, products SET sales.product_price = products.Price WHERE products.id = sales.id_product AND sales.id = ?`
	myInsertSale3 = `UPDATE sales SET income = product_price * quantity`
	mygetAllSales = `SELECT * FROM sales`
	mygetByIdSale = getAllSales + ` WHERE id = $?`
	myupdateSale  = `UPDATE sales SET id_client = ?, id_product = ?, quantity = ?, updated_at = now() WHERE id = ?`
	mydeleteSale  = `DELETE FROM sales WHERE id = ?`
)

type MySqlSale struct {
	db *sql.DB
}

func NewMySqlSale(db *sql.DB) *MySqlSale {
	return &MySqlSale{db}
}

func (p *MySqlSale) Migrate() error {
	db, err := p.db.Prepare(myMigrateSale)
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

func (p *MySqlSale) Insert(m *sales.Model) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	id, err := myfuncInsertSale1(tx, m)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("InsertSale1: %v\n", err)
	}
	err = myfuncInsertSale2(tx, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("InsertSale2: %v\n", err)
	}
	err = myfuncInsertSale3(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("InsertSale3 %v\n", err)
	}
	tx.Commit()
	fmt.Println("Insert completed")
	return nil

}

func myfuncInsertSale1(tx *sql.Tx, m *sales.Model) (uint, error) {
	stmt, err := tx.Prepare(myInsertSale)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(intToNull(int(m.IdClient)), intToNull(int(m.IdProduct)), intToNull(int(m.Quantity)))
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	m.ID = uint(id)
	uid := uint(id)
	fmt.Println("insert number 1 completed")
	return uid, nil
}

func myfuncInsertSale2(tx *sql.Tx, i uint) error {
	stmt, err := tx.Prepare(myInsertSale2)
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

func myfuncInsertSale3(tx *sql.Tx) error {
	stmt, err := tx.Prepare(myInsertSale3)
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

func (p *MySqlSale) GetAll() ([]*sales.Model, error) {
	stmt, err := p.db.Prepare(mygetAllSales)
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

func (p *MySqlSale) GetById(id uint) (*sales.Model, error) {
	stmt, err := p.db.Prepare(mygetByIdSale)
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

func (p *MySqlSale) Update(m *sales.Model) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(myupdateSale)
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
	err = myfuncInsertSale2(tx, m.ID)
	if err != nil {
		return err
	}
	err = myfuncInsertSale3(tx)
	if err != nil {
		return err
	}
	tx.Commit()
	fmt.Println("sale updated succesfuly")
	return nil

}
func (p *MySqlSale) Delete(id uint) error {
	stmt, err := p.db.Prepare(mydeleteSale)
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
