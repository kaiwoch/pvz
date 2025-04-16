package storage

import (
	"database/sql"
	"pvz/internal/storage/migrations/entity"

	"github.com/gofrs/uuid/v5"
)

type ReceptionPostgresStorage struct {
	db *sql.DB
}

func NewReceptionPostgresStorage(db *sql.DB) *ReceptionPostgresStorage {
	return &ReceptionPostgresStorage{db: db}
}

func (r *ReceptionPostgresStorage) GetListReceptionsByPVZId(id uuid.UUID) ([]entity.Receptions, error) {
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
}
