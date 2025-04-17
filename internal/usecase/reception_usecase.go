package usecase

import (
	"database/sql"
	"fmt"
	"pvz/internal/storage"
	"pvz/internal/storage/migrations/entity"

	"github.com/gofrs/uuid/v5"
)

type ReceptionUsecase struct {
	receptionStorage *storage.ReceptionPostgresStorage
}

func NewReceptionUsecase(receptionStorage *storage.ReceptionPostgresStorage) *ReceptionUsecase {
	return &ReceptionUsecase{receptionStorage: receptionStorage}
}

func (r *ReceptionUsecase) CreateReception(id uuid.UUID) (*entity.Receptions, error) {
	_, status, err := r.receptionStorage.GetLastReceptionStatus(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check reception status: %w", err)
	}
	if status == "in_progress" {
		return nil, fmt.Errorf("close previous receipt")
	}
	reception, err := r.receptionStorage.CreateReception(id)
	if err != nil {
		return nil, fmt.Errorf("failed to create new reception: %w", err)
	}
	return reception, nil
}

func (r *ReceptionUsecase) UpdateReceptionStatus(pvz_id uuid.UUID) (*entity.Receptions, error) {
	reception_id, status, err := r.receptionStorage.GetLastReceptionStatus(pvz_id)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to check reception status: %w", err)
	}
	if status == "close" {
		return nil, fmt.Errorf("no available receptions")
	}
	err = r.receptionStorage.UpdateReceptionStatus(reception_id)
	if err != nil {
		return nil, fmt.Errorf("failed to update reception status: %w", err)
	}

	reception, err := r.receptionStorage.GetReceptionById(reception_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get reception: %w", err)
	}

	return reception, nil
}
