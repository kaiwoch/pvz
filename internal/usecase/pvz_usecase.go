package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"pvz/internal/storage"
	"pvz/internal/storage/migrations/entity"
	"time"

	"github.com/gofrs/uuid/v5"
)

type PVZUsecase struct {
	pvzStorage *storage.PVZPostgresStorage
}

type PVZListResponse struct {
	PVZs  []entity.ListPVZ `json:"pvzs"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

func NewPVZUsecase(pvzStorage *storage.PVZPostgresStorage, receptionStorage *storage.ReceptionPostgresStorage, productStorage *storage.ProductPostgresStorage) *PVZUsecase {
	return &PVZUsecase{pvzStorage: pvzStorage}
}

func (p *PVZUsecase) CreatePVZ(id, user_id uuid.UUID, city string, date time.Time) (*entity.PVZ, error) {

	pvz, err := p.pvzStorage.GetPVZById(id)
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("pvz exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check PVZ existence: %w", err)
	}

	pvz, err = p.pvzStorage.CreatePVZ(id, user_id, city, date)
	if err != nil {
		return nil, err
	}
	return pvz, nil
}

func (p *PVZUsecase) GetPVZsWithFilter(ctx context.Context, filter entity.Filter) (*PVZListResponse, error) {
	pvzs, err := p.pvzStorage.GetPVZsWithFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	total, err := p.pvzStorage.CountPVZsWithFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &PVZListResponse{
		PVZs:  pvzs,
		Total: total,
		Page:  filter.Page,
		Limit: filter.Limit,
	}, nil
}
