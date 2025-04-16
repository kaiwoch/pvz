package storage

import (
	"database/sql"
	"pvz/internal/storage/migrations/entity"

	"github.com/gofrs/uuid/v5"
)

type ProductPostgresStorage struct {
	db *sql.DB
}

func NewProductPostgresStorage(db *sql.DB) *ProductPostgresStorage {
	return &ProductPostgresStorage{db: db}
}

func (p *ProductPostgresStorage) GetListProductsByReceptionId(id uuid.UUID) ([]entity.Products, error) {
	query := "SELECT * FROM product WHERE reception_id = $1"
	rows, err := p.db.Query(query, id)

	var productList []entity.Products
	for rows.Next() {
		var product entity.Products
		err = rows.Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionId)
		if err != nil {
			return nil, err
		}
		productList = append(productList, product)
	}
	return productList, nil
}
