package usecase

import (
	"database/sql"
	"fmt"
	"pvz/internal/storage"
	"pvz/internal/storage/migrations/entity"

	"github.com/gofrs/uuid/v5"
)

type ProductUsecase struct {
	productStorage   *storage.ProductPostgresStorage
	receptionStorage *storage.ReceptionPostgresStorage
}

func NewProductUsecase(productStorage *storage.ProductPostgresStorage, receptionStorage *storage.ReceptionPostgresStorage) *ProductUsecase {
	return &ProductUsecase{productStorage: productStorage, receptionStorage: receptionStorage}
}

func (p *ProductUsecase) CreateProduct(id uuid.UUID, product_type string) (*entity.Products, error) {
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
