package storage

import (
	"database/sql"
	"pvz/internal/storage/migrations/entity"
	"time"

	"github.com/gofrs/uuid/v5"
)

type ProductPostgresStorage interface {
	CreateProduct(id uuid.UUID, product_type string) (*entity.Products, error)
	DeleteProduct(product_id uuid.UUID) error
	GetLastProductID(reception_id uuid.UUID) (uuid.UUID, error)
}

type ProductPostgresStorageImpl struct {
	db *sql.DB
}

func NewProductPostgresStorage(db *sql.DB) *ProductPostgresStorageImpl {
	return &ProductPostgresStorageImpl{db: db}
}

func (p *ProductPostgresStorageImpl) CreateProduct(id uuid.UUID, product_type string) (*entity.Products, error) {
	product_id := uuid.Must(uuid.NewV4())
	date := time.Now()
	query := "INSERT INTO product (product_id, date_time, type_name, reception_id) VALUES ($1, $2, $3, $4)"

	_, err := p.db.Exec(query, product_id, date, product_type, id)
	if err != nil {
		return nil, err
	}
	return &entity.Products{ID: product_id, DateTime: date, Type: product_type, ReceptionId: id}, nil
}

func (p *ProductPostgresStorageImpl) DeleteProduct(product_id uuid.UUID) error {
	query := "delete from product where product_id = $1"

	_, err := p.db.Exec(query, product_id)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductPostgresStorageImpl) GetLastProductID(reception_id uuid.UUID) (uuid.UUID, error) {
	var product_id uuid.UUID
	query := "SELECT product_id FROM product WHERE reception_id = $1 ORDER BY date_time DESC LIMIT 1"

	err := p.db.QueryRow(query, reception_id).Scan(&product_id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return product_id, nil
}
