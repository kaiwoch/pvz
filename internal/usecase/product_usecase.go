package usecase

import (
	"database/sql"
	"fmt"
	"pvz/internal/storage"
	"pvz/internal/storage/migrations/entity"

	"github.com/gofrs/uuid/v5"
)

type ProductUsecase interface {
	CreateProduct(id uuid.UUID, product_type string) (*entity.Products, error)
	DeleteLastProduct(pvz_id uuid.UUID) error
}

type ProductUsecaseImpl struct {
	productStorage   storage.ProductPostgresStorage
	receptionStorage storage.ReceptionPostgresStorage
}

func NewProductUsecase(productStorage storage.ProductPostgresStorage, receptionStorage storage.ReceptionPostgresStorage) *ProductUsecaseImpl {
	return &ProductUsecaseImpl{productStorage: productStorage, receptionStorage: receptionStorage}
}

func (p *ProductUsecaseImpl) CreateProduct(id uuid.UUID, product_type string) (*entity.Products, error) {
	reception_id, status, err := p.receptionStorage.GetLastReceptionStatus(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check reception status: %w", err)
	}
	if status == "close" {
		return nil, fmt.Errorf("no available receptions")
	}

	products, err := p.productStorage.CreateProduct(reception_id, product_type)
	if err != nil {
		return nil, fmt.Errorf("failed to create new product: %w", err)
	}
	return products, nil
}

func (p *ProductUsecaseImpl) DeleteLastProduct(pvz_id uuid.UUID) error {
	reception_id, status, err := p.receptionStorage.GetLastReceptionStatus(pvz_id)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check reception status: %w", err)
	}
	if status == "close" {
		return fmt.Errorf("no available receptions")
	}

	product_id, err := p.productStorage.GetLastProductID(reception_id)
	if err != nil {
		return err
	}

	err = p.productStorage.DeleteProduct(product_id)
	if err != nil {
		return err
	}

	return nil
}
