package storage

import (
	"database/sql"
	"pvz/internal/storage/migrations/entity"
	"time"

	"github.com/gofrs/uuid/v5"
)

type PVZPostgresStorage struct {
	db *sql.DB
}

func NewPVZPostgresStorage(db *sql.DB) *PVZPostgresStorage {
	return &PVZPostgresStorage{db: db}
}

func (p *PVZPostgresStorage) CreatePVZ(id, user_id uuid.UUID, city string, date time.Time) (*entity.PVZ, error) {
	query := "INSERT INTO pvz (pvz_id, registration_date, city_name, user_id) VALUES ($1, $2, $3, $4)"
	_, err := p.db.Exec(query, id, date, city, user_id)
	if err != nil {
		return nil, err
	}
	return &entity.PVZ{ID: id, RegistrationDate: date, City: city}, nil // проверить дату
}

func (p *PVZPostgresStorage) GetPVZById(id uuid.UUID) (*entity.PVZ, error) {
	var pvz entity.PVZ
	query := "SELECT * FROM pvz WHERE pvz_id = $1"
	err := p.db.QueryRow(query, id).Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City, &pvz.UserID)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}
	return &pvz, nil
}

/* func (p *PVZPostgresStorage) GetListPVZByUserId(id uuid.UUID) ([]entity.PVZ, error) {
	query := "SELECT * FROM pvz WHERE user_id = $1"
	rows, err := p.db.Query(query, id)

	var pvzList []entity.PVZ
	for rows.Next() {
		var pvz entity.PVZ
		err = rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City)
		if err != nil {
			return nil, err
		}
		pvzList = append(pvzList, pvz)
	}
	return pvzList, nil
} */
