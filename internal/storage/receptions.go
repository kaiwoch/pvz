package storage

import (
	"database/sql"
	"pvz/internal/storage/migrations/entity"
	"time"

	"github.com/gofrs/uuid/v5"
)

type ReceptionPostgresStorage struct {
	db *sql.DB
}

func NewReceptionPostgresStorage(db *sql.DB) *ReceptionPostgresStorage {
	return &ReceptionPostgresStorage{db: db}
}

func (r *ReceptionPostgresStorage) CreateReception(id uuid.UUID) (*entity.Receptions, error) {
	reception_id := uuid.Must(uuid.NewV4())
	date := time.Now()
	status := "in_progress"
	query := "INSERT INTO reception (reception_id, date_time, pvz_id, status_name) VALUES ($1, $2, $3, $4)"

	_, err := r.db.Exec(query, reception_id, date, id, status)
	if err != nil {
		return nil, err
	}

	return &entity.Receptions{ID: reception_id, DateTime: date, PVZID: id, Status: status}, nil
}

func (r *ReceptionPostgresStorage) GetLastReceptionStatus(id uuid.UUID) (uuid.UUID, string, error) {
	var status string
	var reception_id uuid.UUID
	query := "SELECT reception_id, status_name FROM reception WHERE pvz_id = $1 ORDER BY date_time DESC LIMIT 1"

	err := r.db.QueryRow(query, id).Scan(&reception_id, &status)
	if err != nil {
		return uuid.UUID{}, "", err
	}

	return reception_id, status, err
}

func (r *ReceptionPostgresStorage) UpdateReceptionStatus(reception_id uuid.UUID) error {
	query := "UPDATE reception SET status_name = 'close' WHERE reception_id = $1"

	_, err := r.db.Exec(query, reception_id)
	if err != nil {
		return err
	}

	return nil
}

func (r *ReceptionPostgresStorage) GetReceptionById(reception_id uuid.UUID) (*entity.Receptions, error) {
	query := "SELECT * FROM reception WHERE reception_id = $1"
	var date time.Time
	var pvz_id uuid.UUID
	var status string

	err := r.db.QueryRow(query, reception_id).Scan(&reception_id, &date, &pvz_id, &status)
	if err != nil {
		return nil, err
	}
	return &entity.Receptions{ID: reception_id, DateTime: date, PVZID: pvz_id, Status: status}, nil
}
