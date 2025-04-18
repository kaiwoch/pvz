package storage_test

import (
	"pvz/internal/storage"
	"pvz/internal/storage/migrations/entity"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"
)

func TestReceptionPostgresStorage_CreateReception(t *testing.T) {
	pvz_id := uuid.Must(uuid.NewV4())

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := storage.NewReceptionPostgresStorage(db)

	tests := []struct {
		name        string
		id          uuid.UUID
		mock        func()
		expectedErr error
	}{
		{
			name: "success",
			id:   pvz_id,
			mock: func() {
				mock.ExpectExec("INSERT INTO reception").
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), pvz_id, "in_progress").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			reception, err := storage.CreateReception(tt.id)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, reception)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, reception)
				assert.Equal(t, tt.id, reception.PVZID)
				assert.Equal(t, "in_progress", reception.Status)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReceptionPostgresStorage_GetLastReceptionStatus(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := storage.NewReceptionPostgresStorage(db)

	pvz_id := uuid.Must(uuid.NewV4())
	reception_id := uuid.Must(uuid.NewV4())
	status := "close"

	tests := []struct {
		name           string
		input          uuid.UUID
		mock           func()
		expectedID     uuid.UUID
		expectedStatus string
		expectedErr    error
	}{
		{
			name:  "success",
			input: pvz_id,
			mock: func() {
				rows := sqlmock.NewRows([]string{"reception_id", "status_name"}).AddRow(reception_id, status)
				mock.ExpectQuery("SELECT reception_id, status_name FROM reception WHERE pvz_id = \\$1 ORDER BY date_time DESC LIMIT 1").
					WithArgs(pvz_id).WillReturnRows(rows)
			},
			expectedID:     reception_id,
			expectedStatus: status,
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			actualID, actualStatus, err := storage.GetLastReceptionStatus(tt.input)

			assert.Equal(t, tt.expectedID, actualID)
			assert.Equal(t, tt.expectedStatus, actualStatus)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestReceptionPostgresStorage_GetReceptionByID(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	storage := storage.NewReceptionPostgresStorage(db)

	reception_id := uuid.Must(uuid.NewV4())
	date := time.Now()

	tests := []struct {
		name        string
		input       uuid.UUID
		mock        func()
		expected    *entity.Receptions
		expectedErr error
	}{
		{
			name:  "success",
			input: reception_id,
			mock: func() {
				rows := sqlmock.NewRows([]string{"reception_id", "date_time", "pvz_id", "status_name"}).AddRow(reception_id, date, reception_id, "close")
				mock.ExpectQuery("SELECT \\* FROM reception WHERE reception_id = \\$1").
					WithArgs(reception_id).WillReturnRows(rows)
			},
			expected: &entity.Receptions{
				ID:       reception_id,
				PVZID:    reception_id,
				DateTime: date,
				Status:   "close",
			},

			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			reception, err := storage.GetReceptionById(tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
				assert.Nil(t, reception)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.ID, reception.ID)
				assert.Equal(t, tt.expected.PVZID, reception.PVZID)
				assert.Equal(t, tt.expected.Status, reception.Status)
				assert.WithinDuration(t, tt.expected.DateTime, reception.DateTime, time.Second)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
