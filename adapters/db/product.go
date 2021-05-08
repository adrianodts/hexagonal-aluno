package db

import (
	"database/sql"

	"github.com/adrianodts/hexagonal-aluno/application"

	_ "github.com/mattn/go-sqlite3"
)

type ProductDb struct {
	db *sql.DB
}

func NewProductDb(db *sql.DB) *ProductDb {
	return &ProductDb{db: db}
}

// func (p *ProductDb) GetAll() (application.ProductInterface, error) {
// 	var product []application.Product
// 	rows, err := db.Query(`select id, name, price, status from products`)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer rows.Close()
// 	var i = 0
// 	for rows.Next() {
// 		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Status)
// 		if err != nil {
// 			panic(err)
// 		}
// 		product[0] = err
// 		i++
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return &product, nil
// }

func (p *ProductDb) Get(id string) (application.ProductInterface, error) {
	var product application.Product
	stmt, err := p.db.Prepare("select id, name, price, status from products where id=?")
	if err != nil {
		return nil, err
	}
	err = stmt.QueryRow(id).Scan(&product.Id, &product.Name, &product.Price, &product.Status)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *ProductDb) Save(product application.ProductInterface) (application.ProductInterface, error) {
	var rows int
	p.db.QueryRow("select count(*) from products where id=?", product.GetId()).Scan(rows)
	if rows == 0 {
		_, err := p.create(product)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := p.update(product)
		if err != nil {
			return nil, err
		}
	}
	return product, nil
}

func (p *ProductDb) create(product application.ProductInterface) (application.ProductInterface, error) {
	stmt, err := p.db.Prepare(`insert into products (id, name, price, status) values (?,?,?,?)`)
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(
		product.GetId(),
		product.GetName(),
		product.GetPrice(),
		product.GetStatus(),
	)
	if err != nil {
		return nil, err
	}
	err = stmt.Close()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductDb) update(product application.ProductInterface) (application.ProductInterface, error) {
	_, err := p.db.Exec(`update products set name=?, price=?, status=? where id=?`,
		product.GetName(), product.GetPrice(), product.GetStatus(), product.GetId())
	if err != nil {
		return nil, err
	}
	return product, nil
}

// Create(name string, price float64) (ProductInterface, error)
// Enable(product ProductInterface) (ProductInterface, error)
// Disable(product ProductInterface) (ProductInterface, error)
