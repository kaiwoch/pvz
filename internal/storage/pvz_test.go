package storage_test

import (
	"context"
	"database/sql"
	"pvz/internal/storage"
	"pvz/internal/storage/migrations/entity"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

func TestPVZPostgresStorage_CreatePVZ(t *testing.T) {
	pvz_id := uuid.Must(uuid.NewV4())
	user_id := uuid.Must(uuid.NewV4())
	city := "Москва"
	date := time.Now()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := storage.NewPVZPostgresStorage(db)

	tests := []struct {
		name        string
		pvz_id      uuid.UUID
		user_id     uuid.UUID
		city        string
		date        time.Time
		mock        func()
		expectedErr error
	}{
		{
			name:    "success",
			pvz_id:  pvz_id,
			user_id: user_id,
			city:    city,
			date:    date,
			mock: func() {
				mock.ExpectExec("INSERT INTO pvz").
					WithArgs(pvz_id, date, city, user_id).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			pvz, err := storage.CreatePVZ(tt.pvz_id, tt.user_id, tt.city, tt.date)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, pvz)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pvz)
				assert.Equal(t, tt.pvz_id, pvz.ID)
				assert.Equal(t, tt.city, pvz.City)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPVZPostgresStorage_GetPVZById(t *testing.T) {
	pvz_id := uuid.Must(uuid.NewV4())
	user_id := uuid.Must(uuid.NewV4())
	city := "Москва"
	date := time.Now()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := storage.NewPVZPostgresStorage(db)

	tests := []struct {
		name        string
		pvz_id      uuid.UUID
		mock        func()
		expected    *entity.PVZ
		expectedErr error
	}{
		{
			name:   "success",
			pvz_id: pvz_id,
			mock: func() {
				rows := sqlmock.NewRows([]string{"pvz_id", "registration_date", "city_name", "user_id"}).AddRow(pvz_id, date, city, user_id)
				mock.ExpectQuery("SELECT \\* FROM pvz WHERE pvz_id = \\$1").
					WithArgs(pvz_id).WillReturnRows(rows)
			},
			expected: &entity.PVZ{
				ID:               pvz_id,
				RegistrationDate: date,
				City:             city,
				UserID:           user_id,
			},
			expectedErr: nil,
		},
		{
			name:   "not found",
			pvz_id: pvz_id,
			mock: func() {
				mock.ExpectQuery("SELECT \\* FROM pvz WHERE pvz_id = \\$1").
					WithArgs(pvz_id).WillReturnError(sql.ErrNoRows)
			},
			expected:    nil,
			expectedErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			pvz, err := storage.GetPVZById(tt.pvz_id)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, pvz)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pvz)
				assert.Equal(t, tt.expected, pvz)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPVZPostgresStorage_GetPVZsWithFilter(t *testing.T) {
	pvz_id := uuid.Must(uuid.NewV4())
	reception_id := uuid.Must(uuid.NewV4())
	product_id := uuid.Must(uuid.NewV4())
	date := time.Now()

	filter := entity.Filter{
		StartDate: &date,
		EndDate:   &date,
		Page:      1,
		Limit:     10,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := storage.NewPVZPostgresStorage(db)

	tests := []struct {
		name        string
		pvz_id      uuid.UUID
		mock        func()
		expected    []entity.ListPVZ
		expectedErr error
	}{
		{
			name:   "success",
			pvz_id: pvz_id,
			mock: func() {
				rows := sqlmock.NewRows([]string{"pvz_id", "registration_date", "city_name",
					"reception_id", "date_time", "pvz_id", "status_name",
					"product_id", "date_time", "type_name", "reception_id"}).AddRow(pvz_id, date, "Москва",
					reception_id, date, pvz_id, "close",
					product_id, date, "одежда", reception_id)
				mock.ExpectQuery(`
				WITH filtered_receptions AS .*
				SELECT .*
				FROM pvz p
				JOIN filtered_receptions r ON p.pvz_id = r.pvz_id
				LEFT JOIN product pr ON pr.reception_id = r.reception_id
				ORDER BY r.date_time DESC
				LIMIT .* OFFSET .*
			`).WithArgs(filter.StartDate, filter.EndDate, filter.Limit, 0).
					WillReturnRows(rows)
			},
			expected: []entity.ListPVZ{
				{
					Pvz: entity.PVZ{
						ID:               pvz_id,
						RegistrationDate: date,
						City:             "Москва",
					},
					Receptions: []entity.Receptions{
						{
							ID:       reception_id,
							DateTime: date,
							PVZID:    pvz_id,
							Status:   "close",
							Products: []entity.Products{
								{
									ID:          product_id,
									DateTime:    date,
									Type:        "одежда",
									ReceptionId: reception_id,
								},
							},
						},
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			pvz, err := storage.GetPVZsWithFilter(context.Background(), filter)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, pvz)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pvz)
				assert.Equal(t, tt.expected, pvz)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
