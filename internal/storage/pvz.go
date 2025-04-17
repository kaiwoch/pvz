package storage

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *PVZPostgresStorage) GetPVZsWithFilter(ctx context.Context, filter entity.Filter) ([]entity.ListPVZ, error) {
	query := `
		WITH filtered_receptions AS (
			SELECT r.reception_id, r.date_time, r.pvz_id, r.status_name
			FROM reception r
			WHERE ($1::timestamp IS NULL OR r.date_time >= $1)
			AND ($2::timestamp IS NULL OR r.date_time <= $2)
		)
		SELECT 
			p.pvz_id, p.registration_date, p.city_name,
			r.reception_id, r.date_time, r.pvz_id, r.status_name,
			pr.product_id, pr.date_time, pr.type_name, pr.reception_id
		FROM pvz p
		JOIN filtered_receptions r ON p.pvz_id = r.pvz_id
		LEFT JOIN product pr ON pr.reception_id = r.reception_id
		ORDER BY r.date_time DESC
		LIMIT $3 OFFSET $4
	`

	offset := (filter.Page - 1) * filter.Limit

	rows, err := r.db.QueryContext(ctx, query, filter.StartDate, filter.EndDate, filter.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query pvzs: %w", err)
	}
	defer rows.Close()

	var result []entity.ListPVZ
	var currentPVZ *entity.ListPVZ
	var lastPVZID uuid.UUID
	var lastReceptionID uuid.UUID

	for rows.Next() {
		var (
			pvzID               uuid.UUID
			pvzRegistrationDate time.Time
			pvzCity             string
			receptionID         uuid.UUID
			receptionDateTime   time.Time
			receptionPVZID      uuid.UUID
			receptionStatus     string
			productID           uuid.UUID
			productDateTime     time.Time
			productType         string
			productReceptionID  uuid.UUID
		)

		err := rows.Scan(
			&pvzID, &pvzRegistrationDate, &pvzCity,
			&receptionID, &receptionDateTime, &receptionPVZID, &receptionStatus,
			&productID, &productDateTime, &productType, &productReceptionID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if currentPVZ == nil || pvzID != lastPVZID {
			result = append(result, entity.ListPVZ{
				Pvz: entity.PVZ{
					ID:               pvzID,
					RegistrationDate: pvzRegistrationDate,
					City:             pvzCity,
				},
				Receptions: []entity.Receptions{},
			})
			currentPVZ = &result[len(result)-1]
			lastPVZID = pvzID
			lastReceptionID = uuid.Nil
		}

		if receptionID != lastReceptionID {
			currentPVZ.Receptions = append(currentPVZ.Receptions, entity.Receptions{
				ID:       receptionID,
				DateTime: receptionDateTime,
				PVZID:    receptionPVZID,
				Status:   receptionStatus,
				Products: []entity.Products{},
			})
			lastReceptionID = receptionID
		}

		if !productID.IsNil() {
			lastReception := &currentPVZ.Receptions[len(currentPVZ.Receptions)-1]
			lastReception.Products = append(lastReception.Products, entity.Products{
				ID:          productID,
				DateTime:    productDateTime,
				Type:        productType,
				ReceptionId: productReceptionID,
			})
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}

func (r *PVZPostgresStorage) CountPVZsWithFilter(ctx context.Context, filter entity.Filter) (int, error) {
	query := `
		SELECT COUNT(DISTINCT p.pvz_id)
		FROM pvz p
		JOIN reception r ON p.pvz_id = r.pvz_id
		WHERE ($1::timestamp IS NULL OR r.date_time >= $1)
		AND ($2::timestamp IS NULL OR r.date_time <= $2)
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, filter.StartDate, filter.EndDate).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pvzs: %w", err)
	}

	return count, nil
}
