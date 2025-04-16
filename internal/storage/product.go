package storage

import (
	"database/sql"
	"pvz/internal/storage/migrations/entity"
	"time"

	"github.com/gofrs/uuid/v5"
)

type ProductPostgresStorage struct {
	db *sql.DB
}

func NewProductPostgresStorage(db *sql.DB) *ProductPostgresStorage {
	return &ProductPostgresStorage{db: db}
}

/* func (p *ProductPostgresStorage) GetListProductsByReceptionId(id uuid.UUID) ([]entity.Products, error) {
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
} */

func (p *ProductPostgresStorage) CreateProduct(id uuid.UUID, product_type string) (*entity.Products, error) {
	product_id := uuid.Must(uuid.NewV4())
	date := time.Now()
	query := "INSERT INTO product (product_id, date_time, type_name, reception_id) VALUES ($1, $2, $3, $4)"

	_, err := p.db.Exec(query, product_id, date, product_type, id)
	if err != nil {
		return nil, err
	}
	return &entity.Products{ID: product_id, DateTime: date, Type: product_type, ReceptionId: id}, nil
}
