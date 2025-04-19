package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"pvz/internal/storage"
	"pvz/internal/storage/migrations/entity"
	"time"

	"github.com/gofrs/uuid/v5"
)

type PVZUsecase interface {
	CreatePVZ(id, user_id uuid.UUID, city string, date time.Time) (*entity.PVZ, error)
	GetPVZsWithFilter(ctx context.Context, filter entity.Filter) (*PVZListResponse, error)
}

type PVZUsecaseImpl struct {
	pvzStorage storage.PVZPostgresStorage
}

type PVZListResponse struct {
	PVZs  []entity.ListPVZ `json:"pvzs"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

func NewPVZUsecase(pvzStorage storage.PVZPostgresStorage) *PVZUsecaseImpl {
	return &PVZUsecaseImpl{pvzStorage: pvzStorage}
}

func (p *PVZUsecaseImpl) CreatePVZ(id, user_id uuid.UUID, city string, date time.Time) (*entity.PVZ, error) {

	pvz, err := p.pvzStorage.GetPVZById(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check PVZ existence: %w", err)
	}

	if !pvz.ID.IsNil() {
		return nil, errors.New("pvz exists")
	}

	pvz, err = p.pvzStorage.CreatePVZ(id, user_id, city, date)
	if err != nil {
		return nil, err
	}
	return pvz, nil
}

func (p *PVZUsecaseImpl) GetPVZsWithFilter(ctx context.Context, filter entity.Filter) (*PVZListResponse, error) {
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
