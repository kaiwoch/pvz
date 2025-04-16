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

/* func (r *ReceptionPostgresStorage) GetListReceptionsByPVZId(id uuid.UUID) ([]entity.Receptions, error) {
	query := "SELECT * FROM reception WHERE pvz_id = $1"
	rows, err := r.db.Query(query, id)

	var receptionList []entity.Receptions
	for rows.Next() {
		var reception entity.Receptions
		err = rows.Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status)
		if err != nil {
			return nil, err
		}
		receptionList = append(receptionList, reception)
	}
	return receptionList, nil
} */

func (r *ReceptionPostgresStorage) CreateReception(id uuid.UUID) (*entity.Receptions, error) {
	reception_id := uuid.Must(uuid.NewV4())
	date := time.Now()
	status := "in_progress"
	query := "INSERT INTO reception (reception_id, date_time, pvz_id, status_name) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, reception_id, date, id, status)
	if err != nil {
		return nil, err
	}
	return &entity.Receptions{ID: reception_id, DateTime: date, PVZID: id, Status: status, Products: make([]entity.Products, 0)}, nil
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
